package com.hemingway.kayak.ui.spots

import android.os.Bundle
import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup
import android.widget.TextView
import androidx.core.view.isVisible
import androidx.fragment.app.Fragment
import androidx.lifecycle.lifecycleScope
import com.hemingway.kayak.R
import com.hemingway.kayak.data.remote.SpotDto
import com.hemingway.kayak.data.remote.SpotRepository
import com.hemingway.kayak.databinding.FragmentSpotsBinding
import com.hemingway.kayak.location.LocationTracker
import dagger.hilt.android.AndroidEntryPoint
import kotlinx.coroutines.launch
import java.util.Locale
import javax.inject.Inject

/** Spots compartidos por la comunidad: listado y alta. */
@AndroidEntryPoint
class SpotsFragment : Fragment() {

    private var _binding: FragmentSpotsBinding? = null
    private val binding get() = _binding!!

    @Inject
    lateinit var spotRepository: SpotRepository

    @Inject
    lateinit var locationTracker: LocationTracker

    override fun onCreateView(
        inflater: LayoutInflater,
        container: ViewGroup?,
        savedInstanceState: Bundle?
    ): View {
        _binding = FragmentSpotsBinding.inflate(inflater, container, false)
        return binding.root
    }

    override fun onViewCreated(view: View, savedInstanceState: Bundle?) {
        super.onViewCreated(view, savedInstanceState)
        binding.addSpotButton.setOnClickListener { binding.spotForm.isVisible = !binding.spotForm.isVisible }
        binding.useMyLocationButton.setOnClickListener { fillMyLocation() }
        binding.saveSpotButton.setOnClickListener { save() }
        load()
    }

    override fun onResume() {
        super.onResume()
        load()
    }

    private fun load() {
        viewLifecycleOwner.lifecycleScope.launch {
            val spots = spotRepository.all()
            render(spots)
        }
    }

    private fun render(spots: List<SpotDto>) {
        binding.spotsContainer.removeAllViews()
        binding.emptyView.isVisible = spots.isEmpty()
        val inflater = LayoutInflater.from(requireContext())
        spots.forEach { spot ->
            val row = inflater.inflate(R.layout.item_spot, binding.spotsContainer, false)
            row.findViewById<TextView>(R.id.spotName).text = spot.name
            row.findViewById<TextView>(R.id.spotDescription).text =
                listOfNotNull(
                    spot.description?.takeIf { it.isNotBlank() },
                    String.format(Locale.getDefault(), "%.4f, %.4f", spot.latitude, spot.longitude),
                    "★ ${spot.rating ?: 0}"
                ).joinToString(" · ")
            binding.spotsContainer.addView(row)
        }
    }

    private fun fillMyLocation() {
        viewLifecycleOwner.lifecycleScope.launch {
            val location = locationTracker.lastKnownLocation.value
                ?: locationTracker.getLastLocation()
                ?: locationTracker.currentLocation.value
            if (location == null) {
                binding.errorView.text = "Todavía no hay ubicación: activá el tracking primero"
                binding.errorView.isVisible = true
                return@launch
            }
            binding.spotLat.setText(String.format(Locale.US, "%.6f", location.latitude))
            binding.spotLon.setText(String.format(Locale.US, "%.6f", location.longitude))
        }
    }

    private fun save() {
        val name = binding.spotName.text?.toString()?.trim().orEmpty()
        val description = binding.spotDescription.text?.toString()?.trim().orEmpty()
        val lat = binding.spotLat.text?.toString()?.toDoubleOrNull()
        val lon = binding.spotLon.text?.toString()?.toDoubleOrNull()
        if (name.isEmpty() || lat == null || lon == null) {
            binding.errorView.text = "Faltan nombre o coordenadas"
            binding.errorView.isVisible = true
            return
        }
        binding.errorView.isVisible = false
        binding.saveSpotButton.isEnabled = false
        viewLifecycleOwner.lifecycleScope.launch {
            val ok = spotRepository.create(name, description.ifEmpty { null }, lat, lon, 0)
            binding.saveSpotButton.isEnabled = true
            if (!ok) {
                binding.errorView.text = "No se pudo compartir el spot"
                binding.errorView.isVisible = true
                return@launch
            }
            binding.spotName.text?.clear()
            binding.spotDescription.text?.clear()
            binding.spotLat.text?.clear()
            binding.spotLon.text?.clear()
            binding.spotForm.isVisible = false
            load()
        }
    }

    override fun onDestroyView() {
        super.onDestroyView()
        _binding = null
    }
}
