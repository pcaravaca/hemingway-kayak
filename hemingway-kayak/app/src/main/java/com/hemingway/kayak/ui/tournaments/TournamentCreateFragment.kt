package com.hemingway.kayak.ui.tournaments

import android.app.DatePickerDialog
import android.app.TimePickerDialog
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
import com.hemingway.kayak.data.remote.TournamentRepository
import com.hemingway.kayak.databinding.FragmentTournamentCreateBinding
import dagger.hilt.android.AndroidEntryPoint
import kotlinx.coroutines.launch
import java.time.Instant
import java.util.Calendar
import java.util.Date
import java.util.Locale
import javax.inject.Inject

/** Alta de torneos: antes solo se podia unir a uno existente. */
@AndroidEntryPoint
class TournamentCreateFragment : Fragment() {

    private var _binding: FragmentTournamentCreateBinding? = null
    private val binding get() = _binding!!

    @Inject
    lateinit var repository: TournamentRepository

    private var startMillis: Long? = null
    private var endMillis: Long? = null

    override fun onCreateView(
        inflater: LayoutInflater,
        container: ViewGroup?,
        savedInstanceState: Bundle?
    ): View {
        _binding = FragmentTournamentCreateBinding.inflate(inflater, container, false)
        return binding.root
    }

    override fun onViewCreated(view: View, savedInstanceState: Bundle?) {
        super.onViewCreated(view, savedInstanceState)
        binding.tournamentType.adapter = ArrayAdapter(
            requireContext(),
            android.R.layout.simple_spinner_dropdown_item,
            resources.getStringArray(R.array.tournament_types)
        )
        binding.startButton.setOnClickListener { pickDateTime { startMillis = it } }
        binding.endButton.setOnClickListener { pickDateTime { endMillis = it } }
        binding.createButton.setOnClickListener { save() }
        renderDates()
    }

    private fun pickDateTime(onPicked: (Long) -> Unit) {
        val now = Calendar.getInstance()
        DatePickerDialog(
            requireContext(),
            { _, year, month, day ->
                TimePickerDialog(
                    requireContext(),
                    { _, hour, minute ->
                        val cal = Calendar.getInstance().apply {
                            set(year, month, day, hour, minute, 0)
                        }
                        onPicked(cal.timeInMillis)
                        renderDates()
                    },
                    now.get(Calendar.HOUR_OF_DAY),
                    now.get(Calendar.MINUTE),
                    true
                ).show()
            },
            now.get(Calendar.YEAR),
            now.get(Calendar.MONTH),
            now.get(Calendar.DAY_OF_MONTH)
        ).show()
    }

    private fun renderDates() {
        binding.startButton.text = label(startMillis) ?: "Elegir inicio"
        binding.endButton.text = label(endMillis) ?: "Elegir fin"
    }

    private fun label(millis: Long?): String? = millis?.let {
        java.text.SimpleDateFormat("dd/MM/yyyy HH:mm", Locale.getDefault()).format(Date(it))
    }

    private fun save() {
        val name = binding.tournamentName.text?.toString()?.trim().orEmpty()
        if (name.isEmpty()) {
            showError("Ponele un nombre al torneo")
            return
        }
        val start = startMillis
        val end = endMillis
        if (start == null || end == null) {
            showError("Elegí el inicio y el fin")
            return
        }
        if (end <= start) {
            showError("El fin tiene que ser después del inicio")
            return
        }
        binding.errorView.isVisible = false
        binding.createButton.isEnabled = false
        viewLifecycleOwner.lifecycleScope.launch {
            val error = repository.create(
                name = name,
                description = binding.tournamentDescription.text?.toString()?.trim()?.ifEmpty { null },
                type = binding.tournamentType.selectedItem?.toString(),
                maxParticipants = binding.tournamentMax.text?.toString()?.toIntOrNull(),
                startTime = Instant.ofEpochMilli(start).toString(),
                endTime = Instant.ofEpochMilli(end).toString(),
                locationName = binding.tournamentLocation.text?.toString()?.trim()?.ifEmpty { null },
                startLat = binding.tournamentLat.text?.toString()?.toDoubleOrNull() ?: 0.0,
                startLon = binding.tournamentLon.text?.toString()?.toDoubleOrNull() ?: 0.0
            )
            binding.createButton.isEnabled = true
            if (error != null) {
                showError(error)
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
