package com.hemingway.kayak.ui.map

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
import androidx.fragment.app.viewModels
import com.google.android.gms.location.FusedLocationProviderClient
import com.google.android.gms.location.LocationServices
import com.google.android.gms.location.Priority
import com.google.android.gms.tasks.CancellationTokenSource
import com.hemingway.kayak.R
import com.hemingway.kayak.data.FishingSpot
import com.hemingway.kayak.databinding.FragmentMapBinding
import dagger.hilt.android.AndroidEntryPoint
import org.osmdroid.config.Configuration
import org.osmdroid.tileprovider.tilesource.TileSourceFactory
import org.osmdroid.util.GeoPoint
import org.osmdroid.views.MapView
import org.osmdroid.views.overlay.Marker
import java.util.Date

@AndroidEntryPoint
class MapFragment : Fragment() {

    private var _binding: FragmentMapBinding? = null
    private val binding get() = _binding!!

    private val viewModel: MapViewModel by viewModels()

    private lateinit var mapView: MapView
    private lateinit var fusedLocationClient: FusedLocationProviderClient

    private val markers = mutableMapOf<Int, Marker>()
    private var currentLocationMarker: Marker? = null

    private val locationPermissionRequest = registerForActivityResult(
        ActivityResultContracts.RequestMultiplePermissions()
    ) { permissions ->
        when {
            permissions[Manifest.permission.ACCESS_FINE_LOCATION] == true -> {
                enableMyLocation()
            }
            else -> {
            }
        }
    }

    override fun onCreateView(
        inflater: LayoutInflater,
        container: ViewGroup?,
        savedInstanceState: Bundle?
    ): View {
        _binding = FragmentMapBinding.inflate(inflater, container, false)
        return binding.root
    }

    override fun onViewCreated(view: View, savedInstanceState: Bundle?) {
        super.onViewCreated(view, savedInstanceState)

        Configuration.getInstance().userAgentValue = requireContext().packageName

        mapView = binding.map
        mapView.setTileSource(TileSourceFactory.MAPNIK)
        mapView.setMultiTouchControls(true)
        mapView.controller.setZoom(15.0)

        fusedLocationClient = LocationServices.getFusedLocationProviderClient(requireActivity())

        setupMapListeners()
        setupClickListeners()
        observeViewModel()
        checkLocationPermissions()
    }

    private fun setupMapListeners() {
        mapView.setOnLongClickListener { event ->
            val projection = mapView.projection
            val geoPoint = projection.fromPixels(event.x.toInt(), event.y.toInt()) as? GeoPoint
            geoPoint?.let { showAddFishingSpotDialog(it) }
            true
        }
    }

    private fun setupClickListeners() {
        binding.fabAddFishingSpot.setOnClickListener {
            addCurrentLocationAsFishingSpot()
        }

        binding.btnCloseInfo.setOnClickListener {
            binding.cardFishingSpotInfo.isVisible = false
        }

        binding.btnEditFishingSpot.setOnClickListener {
        }

        binding.btnDeleteFishingSpot.setOnClickListener {
            viewModel.selectedSpot.value?.let { spot ->
                viewModel.deleteFishingSpot(spot)
                binding.cardFishingSpotInfo.isVisible = false
            }
        }
    }

    private fun observeViewModel() {
        viewModel.fishingSpots.observe(viewLifecycleOwner) { spots ->
            updateMapMarkers(spots)
        }

        viewModel.selectedSpot.observe(viewLifecycleOwner) { spot ->
            spot?.let { showSpotInfo(it) }
        }
    }

    private fun checkLocationPermissions() {
        when {
            ContextCompat.checkSelfPermission(
                requireContext(),
                Manifest.permission.ACCESS_FINE_LOCATION
            ) == PackageManager.PERMISSION_GRANTED -> {
                enableMyLocation()
            }
            else -> {
                locationPermissionRequest.launch(
                    arrayOf(
                        Manifest.permission.ACCESS_FINE_LOCATION,
                        Manifest.permission.ACCESS_COARSE_LOCATION
                    )
                )
            }
        }
    }

    private fun enableMyLocation() {
        if (ContextCompat.checkSelfPermission(
                requireContext(),
                Manifest.permission.ACCESS_FINE_LOCATION
            ) != PackageManager.PERMISSION_GRANTED
        ) {
            return
        }

        val cancellationToken = CancellationTokenSource()
        fusedLocationClient.getCurrentLocation(
            Priority.PRIORITY_HIGH_ACCURACY,
            cancellationToken.token
        ).addOnSuccessListener { location ->
            location?.let {
                val geoPoint = GeoPoint(it.latitude, it.longitude)
                mapView.controller.animateTo(geoPoint)

                if (currentLocationMarker == null) {
                    currentLocationMarker = Marker(mapView).apply {
                        position = geoPoint
                        setAnchor(Marker.ANCHOR_CENTER, Marker.ANCHOR_BOTTOM)
                        title = getString(R.string.location_permission_required)
                    }
                    mapView.overlays.add(currentLocationMarker)
                } else {
                    currentLocationMarker?.position = geoPoint
                }
                mapView.invalidate()

                viewModel.setCurrentLocation(geoPoint)
            }
        }
    }

    private fun updateMapMarkers(fishingSpots: List<FishingSpot>) {
        val currentIds = fishingSpots.map { it.id }
        val idsToRemove = markers.keys.filter { it !in currentIds }
        idsToRemove.forEach { id ->
            markers[id]?.let { mapView.overlays.remove(it) }
            markers.remove(id)
        }

        for (spot in fishingSpots) {
            val position = GeoPoint(spot.latitude, spot.longitude)

            if (markers.containsKey(spot.id)) {
                markers[spot.id]?.position = position
                markers[spot.id]?.title = spot.name
                markers[spot.id]?.snippet = spot.description
            } else {
                val marker = Marker(mapView).apply {
                    this.position = position
                    setAnchor(Marker.ANCHOR_CENTER, Marker.ANCHOR_BOTTOM)
                    title = spot.name
                    snippet = spot.description
                    id = spot.id.toString()
                    setOnMarkerClickListener { clickedMarker, _ ->
                        viewModel.selectFishingSpot(spot.id)
                        true
                    }
                }
                markers[spot.id] = marker
                mapView.overlays.add(marker)
            }
        }
        mapView.invalidate()
    }

    private fun showSpotInfo(spot: FishingSpot) {
        binding.cardFishingSpotInfo.isVisible = true
        binding.tvFishingSpotName.text = spot.name
        binding.tvFishingSpotDescription.text = spot.description
        binding.ratingFishingSpot.rating = spot.rating.toFloat()
        binding.tvFishingSpotFishCaught.text = getString(R.string.fishing_spot_fish_caught, spot.fishCaught)
    }

    private fun addCurrentLocationAsFishingSpot() {
        if (ContextCompat.checkSelfPermission(
                requireContext(),
                Manifest.permission.ACCESS_FINE_LOCATION
            ) != PackageManager.PERMISSION_GRANTED
        ) {
            return
        }

        val cancellationToken = CancellationTokenSource()
        fusedLocationClient.getCurrentLocation(
            Priority.PRIORITY_HIGH_ACCURACY,
            cancellationToken.token
        ).addOnSuccessListener { location ->
            location?.let {
                val geoPoint = GeoPoint(it.latitude, it.longitude)
                showAddFishingSpotDialog(geoPoint)
            }
        }
    }

    private fun showAddFishingSpotDialog(geoPoint: GeoPoint) {
        val newSpot = FishingSpot(
            name = getString(R.string.add_fishing_spot),
            description = "Lat: ${geoPoint.latitude}, Lon: ${geoPoint.longitude}",
            latitude = geoPoint.latitude,
            longitude = geoPoint.longitude,
            createdAt = Date()
        )
        viewModel.saveFishingSpot(newSpot)
    }

    override fun onResume() {
        super.onResume()
        mapView.onResume()
    }

    override fun onPause() {
        super.onPause()
        mapView.onPause()
    }

    override fun onDestroyView() {
        super.onDestroyView()
        _binding = null
    }
}
