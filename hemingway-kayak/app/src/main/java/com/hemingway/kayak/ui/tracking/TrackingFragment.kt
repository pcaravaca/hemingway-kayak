package com.hemingway.kayak.ui.tracking

import android.Manifest
import android.content.pm.PackageManager
import android.os.Bundle
import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup
import androidx.activity.result.contract.ActivityResultContracts
import androidx.core.content.ContextCompat
import androidx.core.view.isVisible
import androidx.fragment.app.Fragment
import androidx.lifecycle.Lifecycle
import androidx.lifecycle.lifecycleScope
import androidx.lifecycle.repeatOnLifecycle
import com.hemingway.kayak.data.remote.TrackingUploader
import com.hemingway.kayak.databinding.FragmentTrackingBinding
import dagger.hilt.android.AndroidEntryPoint
import kotlinx.coroutines.launch
import org.osmdroid.config.Configuration
import org.osmdroid.tileprovider.tilesource.TileSourceFactory
import org.osmdroid.util.GeoPoint
import org.osmdroid.views.overlay.Marker
import org.osmdroid.views.overlay.Polyline
import javax.inject.Inject

@AndroidEntryPoint
class TrackingFragment : Fragment() {

    private var _binding: FragmentTrackingBinding? = null
    private val binding get() = _binding!!

    @Inject
    lateinit var trackingUploader: TrackingUploader

    private lateinit var mapView: org.osmdroid.views.MapView
    private val route = Polyline()
    private var lastMarker: Marker? = null
    private var pointCount = 0

    private val permissionLauncher = registerForActivityResult(
        ActivityResultContracts.RequestMultiplePermissions()
    ) { result ->
        if (result[Manifest.permission.ACCESS_FINE_LOCATION] == true ||
            result[Manifest.permission.ACCESS_COARSE_LOCATION] == true
        ) {
            start()
        }
    }

    override fun onCreateView(
        inflater: LayoutInflater,
        container: ViewGroup?,
        savedInstanceState: Bundle?
    ): View {
        _binding = FragmentTrackingBinding.inflate(inflater, container, false)
        return binding.root
    }

    override fun onViewCreated(view: View, savedInstanceState: Bundle?) {
        super.onViewCreated(view, savedInstanceState)

        Configuration.getInstance().userAgentValue = requireContext().packageName

        mapView = binding.map
        mapView.setTileSource(TileSourceFactory.MAPNIK)
        mapView.setMultiTouchControls(true)
        mapView.setUseDataConnection(true)
        mapView.overlays.add(route)
        mapView.controller.setZoom(12.0)
        mapView.controller.setCenter(GeoPoint(9.9382, -84.8371))

        binding.toggleButton.setOnClickListener {
            if (trackingUploader.isTracking.value) {
                trackingUploader.stop()
            } else {
                requestOrStart()
            }
        }

        viewLifecycleOwner.lifecycleScope.launch {
            viewLifecycleOwner.repeatOnLifecycle(Lifecycle.State.STARTED) {
                trackingUploader.isTracking.collect { tracking ->
                    renderState(tracking)
                }
            }
        }
    }

    private fun requestOrStart() {
        val fine = ContextCompat.checkSelfPermission(requireContext(), Manifest.permission.ACCESS_FINE_LOCATION)
        val coarse = ContextCompat.checkSelfPermission(requireContext(), Manifest.permission.ACCESS_COARSE_LOCATION)
        if (fine == PackageManager.PERMISSION_GRANTED || coarse == PackageManager.PERMISSION_GRANTED) {
            start()
        } else {
            permissionLauncher.launch(
                arrayOf(
                    Manifest.permission.ACCESS_FINE_LOCATION,
                    Manifest.permission.ACCESS_COARSE_LOCATION
                )
            )
        }
    }

    private fun start() {
        trackingUploader.start()
        observeLocations()
    }

    private fun observeLocations() {
        viewLifecycleOwner.lifecycleScope.launch {
            trackingUploader.tracker.currentLocation.collect { location ->
                if (location == null) return@collect
                if (!trackingUploader.isTracking.value) return@collect
                drawPoint(location.latitude, location.longitude)
            }
        }
    }

    private fun drawPoint(latitude: Double, longitude: Double) {
        val point = GeoPoint(latitude, longitude)
        route.addPoint(point)
        lastMarker?.let { mapView.overlays.remove(it) }
        lastMarker = Marker(mapView).apply {
            position = point
            setAnchor(Marker.ANCHOR_CENTER, Marker.ANCHOR_BOTTOM)
        }
        lastMarker?.let { mapView.overlays.add(it) }
        mapView.controller.setCenter(point)
        mapView.invalidate()

        pointCount += 1
        binding.pointsText.text = "$pointCount puntos"
    }

    private fun renderState(tracking: Boolean) {
        binding.toggleButton.text = if (tracking) "Detener tracking" else "Iniciar tracking"
        binding.statusText.text = if (tracking) {
            "Tracking activo · enviando al server"
        } else {
            "Sin tracking"
        }
    }

    override fun onDestroyView() {
        super.onDestroyView()
        _binding = null
    }
}
