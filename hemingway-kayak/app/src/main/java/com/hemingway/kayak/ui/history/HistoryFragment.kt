package com.hemingway.kayak.ui.history

import android.os.Bundle
import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup
import androidx.core.view.isVisible
import androidx.fragment.app.Fragment
import androidx.lifecycle.lifecycleScope
import com.hemingway.kayak.data.remote.TripDto
import com.hemingway.kayak.data.remote.UserRepository
import com.hemingway.kayak.databinding.FragmentHistoryBinding
import com.hemingway.kayak.databinding.ItemTripBinding
import dagger.hilt.android.AndroidEntryPoint
import kotlinx.coroutines.launch
import javax.inject.Inject

@AndroidEntryPoint
class HistoryFragment : Fragment() {

    private var _binding: FragmentHistoryBinding? = null
    private val binding get() = _binding!!

    @Inject
    lateinit var userRepository: UserRepository

    override fun onCreateView(
        inflater: LayoutInflater,
        container: ViewGroup?,
        savedInstanceState: Bundle?
    ): View {
        _binding = FragmentHistoryBinding.inflate(inflater, container, false)
        return binding.root
    }

    override fun onViewCreated(view: View, savedInstanceState: Bundle?) {
        super.onViewCreated(view, savedInstanceState)
        load()
    }

    private fun load() {
        binding.progressBar.isVisible = true
        binding.errorView.isVisible = false
        viewLifecycleOwner.lifecycleScope.launch {
            val trips = userRepository.myTrips()
            binding.progressBar.isVisible = false
            if (trips.isEmpty()) {
                binding.emptyView.isVisible = true
            } else {
                render(trips)
            }
        }
    }

    private fun render(trips: List<TripDto>) {
        val container = binding.tripContainer
        container.removeAllViews()
        val inflater = LayoutInflater.from(requireContext())
        for (trip in trips) {
            val row = ItemTripBinding.inflate(inflater, container, false)
            row.titleText.text = trip.title ?: "Salida sin título"
            row.statusText.text = trip.status ?: ""
            row.metaText.text = listOfNotNull(
                trip.type,
                trip.startTime?.take(10),
                trip.distance?.let { if (it > 0) "$it km" else null }
            ).joinToString(" · ")
            row.descriptionText.text = trip.description.orEmpty()
            container.addView(row.root)
        }
    }

    override fun onDestroyView() {
        super.onDestroyView()
        _binding = null
    }
}
