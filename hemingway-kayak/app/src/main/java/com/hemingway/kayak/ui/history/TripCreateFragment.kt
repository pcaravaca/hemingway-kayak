package com.hemingway.kayak.ui.history

import android.os.Bundle
import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup
import android.widget.ArrayAdapter
import androidx.core.view.isVisible
import androidx.fragment.app.Fragment
import androidx.lifecycle.lifecycleScope
import androidx.navigation.fragment.findNavController
import com.hemingway.kayak.R
import com.hemingway.kayak.data.remote.TripRepository
import com.hemingway.kayak.databinding.FragmentTripCreateBinding
import dagger.hilt.android.AndroidEntryPoint
import kotlinx.coroutines.launch
import javax.inject.Inject

/** Alta manual de salidas: la app ya las creaba con el tracking, pero no registrarlas a mano. */
@AndroidEntryPoint
class TripCreateFragment : Fragment() {

    private var _binding: FragmentTripCreateBinding? = null
    private val binding get() = _binding!!

    @Inject
    lateinit var repository: TripRepository

    override fun onCreateView(
        inflater: LayoutInflater,
        container: ViewGroup?,
        savedInstanceState: Bundle?
    ): View {
        _binding = FragmentTripCreateBinding.inflate(inflater, container, false)
        return binding.root
    }

    override fun onViewCreated(view: View, savedInstanceState: Bundle?) {
        super.onViewCreated(view, savedInstanceState)
        binding.tripType.adapter = ArrayAdapter(
            requireContext(), android.R.layout.simple_spinner_dropdown_item,
            resources.getStringArray(R.array.trip_types)
        )
        binding.createButton.setOnClickListener { save() }
    }

    private fun save() {
        val title = binding.tripTitle.text?.toString()?.trim().orEmpty()
        if (title.isEmpty()) {
            showError("La salida necesita un título")
            return
        }
        binding.errorView.isVisible = false
        binding.createButton.isEnabled = false
        viewLifecycleOwner.lifecycleScope.launch {
            val id = repository.create(
                title = title,
                description = binding.tripDescription.text?.toString()?.trim()?.ifEmpty { null },
                type = binding.tripType.selectedItem?.toString(),
                startLat = binding.tripLat.text?.toString()?.toDoubleOrNull() ?: 0.0,
                startLon = binding.tripLon.text?.toString()?.toDoubleOrNull() ?: 0.0
            )
            binding.createButton.isEnabled = true
            if (id == null) {
                showError("No se pudo guardar la salida")
                return@launch
            }
            findNavController().navigateUp()
        }
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
