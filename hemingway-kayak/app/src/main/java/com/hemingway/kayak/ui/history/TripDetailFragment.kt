package com.hemingway.kayak.ui.history

import android.os.Bundle
import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup
import android.widget.ArrayAdapter
import android.widget.TextView
import androidx.core.view.isVisible
import androidx.fragment.app.Fragment
import androidx.lifecycle.lifecycleScope
import androidx.navigation.fragment.findNavController
import com.google.android.material.dialog.MaterialAlertDialogBuilder
import com.hemingway.kayak.R
import com.hemingway.kayak.data.remote.TripDto
import com.hemingway.kayak.data.remote.TripPointDto
import com.hemingway.kayak.data.remote.TripRepository
import com.hemingway.kayak.databinding.FragmentTripDetailBinding
import dagger.hilt.android.AndroidEntryPoint
import kotlinx.coroutines.launch
import java.util.Locale
import javax.inject.Inject

/** Detalle y edicion de una salida: antes la app solo la listaba. */
@AndroidEntryPoint
class TripDetailFragment : Fragment() {

    private var _binding: FragmentTripDetailBinding? = null
    private val binding get() = _binding!!

    @Inject
    lateinit var repository: TripRepository

    private var tripId: String? = null
    private var trip: TripDto? = null

    override fun onCreateView(
        inflater: LayoutInflater,
        container: ViewGroup?,
        savedInstanceState: Bundle?
    ): View {
        _binding = FragmentTripDetailBinding.inflate(inflater, container, false)
        return binding.root
    }

    override fun onViewCreated(view: View, savedInstanceState: Bundle?) {
        super.onViewCreated(view, savedInstanceState)
        tripId = arguments?.getString("id")
        binding.tripType.adapter = ArrayAdapter(
            requireContext(), android.R.layout.simple_spinner_dropdown_item,
            resources.getStringArray(R.array.trip_types)
        )
        binding.tripStatus.adapter = ArrayAdapter(
            requireContext(), android.R.layout.simple_spinner_dropdown_item,
            resources.getStringArray(R.array.trip_statuses)
        )
        binding.saveButton.setOnClickListener { save() }
        binding.deleteButton.setOnClickListener { confirmDelete() }
        binding.addPointButton.setOnClickListener { addPoint() }
        load()
    }

    override fun onResume() {
        super.onResume()
        if (tripId != null) load()
    }

    private fun load() {
        val id = tripId ?: return
        viewLifecycleOwner.lifecycleScope.launch {
            val loaded = repository.trip(id)
            if (loaded == null) {
                binding.errorView.text = "No se pudo cargar la salida"
                binding.errorView.isVisible = true
                return@launch
            }
            trip = loaded
            render(loaded)
            renderPoints(repository.points(id))
        }
    }

    private fun render(t: TripDto) {
        binding.statusText.text = t.status ?: ""
        binding.tripTitle.setText(t.title.orEmpty())
        binding.tripDescription.setText(t.description.orEmpty())
        select(binding.tripType, R.array.trip_types, t.type)
        select(binding.tripStatus, R.array.trip_statuses, t.status)
        binding.metaText.text = listOfNotNull(
            t.startTime?.take(10),
            t.distance?.takeIf { it > 0 }?.let { String.format(Locale.getDefault(), "%.2f km", it) }
        ).joinToString(" · ")
    }

    private fun select(spinner: android.widget.Spinner, arrayRes: Int, value: String?) {
        val options = resources.getStringArray(arrayRes)
        val index = options.indexOf(value)
        if (index >= 0) spinner.setSelection(index)
    }

    private fun renderPoints(points: List<TripPointDto>) {
        binding.pointsSummary.text = "${points.size} puntos"
        binding.pointsContainer.removeAllViews()
        val inflater = LayoutInflater.from(requireContext())
        points.forEach { p ->
            val row = inflater.inflate(android.R.layout.simple_list_item_1, binding.pointsContainer, false)
            (row as TextView).text = String.format(
                Locale.getDefault(), "%.5f, %.5f", p.latitude, p.longitude
            )
            (row as TextView).setTextColor(resources.getColor(R.color.logo_turquoise, null))
            binding.pointsContainer.addView(row)
        }
    }

    private fun save() {
        val id = tripId ?: return
        val title = binding.tripTitle.text?.toString()?.trim().orEmpty()
        if (title.isEmpty()) {
            showError("La salida necesita un título")
            return
        }
        binding.errorView.isVisible = false
        binding.saveButton.isEnabled = false
        viewLifecycleOwner.lifecycleScope.launch {
            val error = repository.update(
                id,
                title,
                binding.tripDescription.text?.toString()?.trim()?.ifEmpty { null },
                binding.tripType.selectedItem?.toString(),
                binding.tripStatus.selectedItem?.toString()
            )
            binding.saveButton.isEnabled = true
            if (error != null) {
                showError(error)
                return@launch
            }
            load()
        }
    }

    private fun addPoint() {
        val id = tripId ?: return
        val lat = binding.pointLat.text?.toString()?.toDoubleOrNull()
        val lon = binding.pointLon.text?.toString()?.toDoubleOrNull()
        if (lat == null || lon == null) {
            showError("Revisá las coordenadas del punto")
            return
        }
        binding.errorView.isVisible = false
        binding.addPointButton.isEnabled = false
        viewLifecycleOwner.lifecycleScope.launch {
            val error = repository.addPoint(id, lat, lon)
            binding.addPointButton.isEnabled = true
            if (error != null) {
                showError(error)
                return@launch
            }
            renderPoints(repository.points(id))
        }
    }

    private fun confirmDelete() {
        val id = tripId ?: return
        MaterialAlertDialogBuilder(requireContext())
            .setTitle("Eliminar salida")
            .setMessage("Se borra la salida y todos sus puntos. No se puede deshacer.")
            .setNegativeButton("Cancelar", null)
            .setPositiveButton("Eliminar") { _, _ ->
                viewLifecycleOwner.lifecycleScope.launch {
                    val error = repository.delete(id)
                    if (error != null) {
                        showError(error)
                        return@launch
                    }
                    findNavController().navigateUp()
                }
            }
            .show()
    }

    private fun showError(message: String) {
        binding.errorView.text = message
        binding.errorView.isVisible = true
    }

    override fun onDestroyView() {
        super.onDestroyView()
        _binding = null
    }
}
