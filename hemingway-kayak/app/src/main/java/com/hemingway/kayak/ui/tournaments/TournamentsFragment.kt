package com.hemingway.kayak.ui.tournaments

import android.os.Bundle
import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup
import androidx.core.view.isVisible
import androidx.fragment.app.Fragment
import androidx.lifecycle.lifecycleScope
import androidx.navigation.fragment.findNavController
import com.hemingway.kayak.R
import com.hemingway.kayak.data.remote.TournamentDto
import com.hemingway.kayak.data.remote.TournamentRepository
import com.hemingway.kayak.databinding.FragmentTournamentsBinding
import com.hemingway.kayak.databinding.ItemTournamentBinding
import dagger.hilt.android.AndroidEntryPoint
import kotlinx.coroutines.launch
import javax.inject.Inject

@AndroidEntryPoint
class TournamentsFragment : Fragment() {

    private var _binding: FragmentTournamentsBinding? = null
    private val binding get() = _binding!!

    @Inject
    lateinit var repository: TournamentRepository

    override fun onCreateView(
        inflater: LayoutInflater,
        container: ViewGroup?,
        savedInstanceState: Bundle?
    ): View {
        _binding = FragmentTournamentsBinding.inflate(inflater, container, false)
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
            val tournaments = repository.tournaments()
            binding.progressBar.isVisible = false
            if (tournaments.isEmpty()) {
                binding.emptyView.isVisible = true
            } else {
                render(tournaments)
            }
        }
    }

    private fun render(tournaments: List<TournamentDto>) {
        val container = binding.tournamentContainer
        container.removeAllViews()
        val inflater = LayoutInflater.from(requireContext())

        for (tournament in tournaments) {
            val row = ItemTournamentBinding.inflate(inflater, container, false)
            row.nameText.text = tournament.name
            row.statusText.text = tournament.status ?: ""
            row.metaText.text = listOfNotNull(
                tournament.type,
                tournament.participantsCount?.let { "$it participantes" }
            ).joinToString(" · ")

            row.joinButton.setOnClickListener { join(tournament, row.joinButton) }
            row.openButton.setOnClickListener {
                findNavController().navigate(
                    R.id.action_tournaments_to_detail,
                    Bundle().apply { putString("id", tournament.id) }
                )
            }
            container.addView(row.root)
        }
    }

    private fun join(tournament: TournamentDto, button: View) {
        button.isEnabled = false
        viewLifecycleOwner.lifecycleScope.launch {
            val error = repository.join(tournament.id, null)
            if (error != null) {
                binding.errorView.text = error
                binding.errorView.isVisible = true
            }
            load()
        }
    }

    override fun onDestroyView() {
        super.onDestroyView()
        _binding = null
    }
}
