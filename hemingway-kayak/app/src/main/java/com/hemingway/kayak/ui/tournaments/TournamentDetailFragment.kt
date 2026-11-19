package com.hemingway.kayak.ui.tournaments

import android.net.Uri
import android.os.Bundle
import androidx.activity.result.contract.ActivityResultContracts
import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup
import androidx.core.view.isVisible
import androidx.fragment.app.Fragment
import androidx.lifecycle.lifecycleScope
import com.hemingway.kayak.data.remote.CatchDto
import com.hemingway.kayak.data.remote.ParticipantDto
import com.hemingway.kayak.data.remote.RankingDto
import com.hemingway.kayak.data.remote.TournamentDto
import com.hemingway.kayak.data.remote.TournamentRepository
import com.hemingway.kayak.data.remote.PhotoRepository
import com.hemingway.kayak.data.remote.UserRepository
import com.hemingway.kayak.databinding.FragmentTournamentDetailBinding
import com.hemingway.kayak.databinding.ItemCatchBinding
import com.hemingway.kayak.databinding.ItemRankingBinding
import dagger.hilt.android.AndroidEntryPoint
import kotlinx.coroutines.launch
import javax.inject.Inject

@AndroidEntryPoint
class TournamentDetailFragment : Fragment() {

    private var _binding: FragmentTournamentDetailBinding? = null
    private val binding get() = _binding!!

    @Inject
    lateinit var repository: TournamentRepository

    @Inject
    lateinit var userRepository: UserRepository

    @Inject
    lateinit var photoRepository: PhotoRepository

    private val tournamentId: String by lazy { arguments?.getString("id").orEmpty() }
    private var myParticipant: ParticipantDto? = null
    private var selectedPhoto: Uri? = null

    override fun onCreateView(
        inflater: LayoutInflater,
        container: ViewGroup?,
        savedInstanceState: Bundle?
    ): View {
        _binding = FragmentTournamentDetailBinding.inflate(inflater, container, false)
        return binding.root
    }

    override fun onViewCreated(view: View, savedInstanceState: Bundle?) {
        super.onViewCreated(view, savedInstanceState)
        binding.joinButton.setOnClickListener { join() }
        binding.catchButton.setOnClickListener { recordCatch() }
        binding.photoButton.setOnClickListener { pickPhoto() }
        load()
    }

    private val pickImage = registerForActivityResult(
        ActivityResultContracts.GetContent()
    ) { uri ->
        if (uri != null) {
            selectedPhoto = uri
            binding.photoButton.text = "Foto elegida ✓"
        }
    }

    private fun pickPhoto() {
        pickImage.launch("image/*")
    }

    private fun load() {
        if (tournamentId.isEmpty()) {
            binding.errorView.text = "Torneo inválido"
            binding.errorView.isVisible = true
            return
        }
        binding.progressBar.isVisible = true
        binding.errorView.isVisible = false
        viewLifecycleOwner.lifecycleScope.launch {
            val me = userRepository.me()
            val participants = repository.participants(tournamentId)
            myParticipant = me?.let { current -> participants.firstOrNull { it.userId == current.id } }

            val tournament = repository.tournament(tournamentId)
            binding.progressBar.isVisible = false
            if (tournament == null) {
                binding.errorView.text = "No se pudo cargar el torneo"
                binding.errorView.isVisible = true
                return@launch
            }
            renderHeader(tournament)
            renderRanking(repository.ranking(tournamentId))
            renderCatches(repository.catches(tournamentId))
            updateJoinState()
        }
    }

    private fun renderHeader(tournament: TournamentDto) {
        binding.nameText.text = tournament.name
        binding.metaText.text = listOfNotNull(
            tournament.status,
            tournament.type,
            tournament.locationName,
            tournament.startTime?.take(10)
        ).joinToString(" · ")
        binding.descriptionText.text = tournament.description.orEmpty()
    }

    private fun renderRanking(ranking: List<RankingDto>) {
        val container = binding.rankingContainer
        container.removeAllViews()
        val inflater = LayoutInflater.from(requireContext())
        for (row in ranking) {
            val item = ItemRankingBinding.inflate(inflater, container, false)
            item.rankText.text = (row.rank ?: 0).toString()
            item.nameText.text = row.displayName ?: row.username ?: "Kayakista"
            item.subText.text = "${row.catchesCount ?: 0} capturas"
            item.scoreText.text = (row.score ?: 0.0).let { if (it % 1.0 == 0.0) it.toInt().toString() else it.toString() }
            container.addView(item.root)
        }
    }

    private fun renderCatches(catches: List<CatchDto>) {
        val container = binding.catchesContainer
        container.removeAllViews()
        val inflater = LayoutInflater.from(requireContext())
        for (catch in catches) {
            val item = ItemCatchBinding.inflate(inflater, container, false)
            item.speciesText.text = catch.species ?: "Sin especie"
            item.weightText.text = "${catch.weight ?: 0.0} kg"
            item.whoText.text = catch.username.orEmpty()
            container.addView(item.root)
        }
    }

    private fun updateJoinState() {
        val joined = myParticipant != null
        binding.joinButton.isEnabled = !joined
        binding.joinButton.text = if (joined) "Ya estás participando" else "Unirme al torneo"
        binding.catchButton.isEnabled = joined
    }

    private fun join() {
        binding.joinButton.isEnabled = false
        viewLifecycleOwner.lifecycleScope.launch {
            val error = repository.join(tournamentId, null)
            if (error != null) {
                binding.errorView.text = error
                binding.errorView.isVisible = true
                updateJoinState()
                return@launch
            }
            load()
        }
    }

    private fun recordCatch() {
        val participant = myParticipant
        if (participant == null) {
            binding.errorView.text = "Unite al torneo antes de cargar capturas"
            binding.errorView.isVisible = true
            return
        }
        val species = binding.speciesInput.text?.toString()?.trim().orEmpty()
        val weight = binding.weightInput.text?.toString()?.trim()?.toDoubleOrNull()
        if (species.isEmpty() || weight == null || weight <= 0.0) {
            binding.errorView.text = "Especie y peso son obligatorios"
            binding.errorView.isVisible = true
            return
        }
        binding.catchButton.isEnabled = false
        viewLifecycleOwner.lifecycleScope.launch {
            val photoUrl = selectedPhoto?.let { photoRepository.upload(it) }
            val error = repository.recordCatch(
                tournamentId = tournamentId,
                participantId = participant.id,
                species = species,
                weight = weight,
                latitude = null,
                longitude = null,
                photoUrl = photoUrl
            )
            binding.catchButton.isEnabled = true
            if (error != null) {
                binding.errorView.text = error
                binding.errorView.isVisible = true
                return@launch
            }
            binding.speciesInput.text?.clear()
            binding.weightInput.text?.clear()
            selectedPhoto = null
            binding.photoButton.text = "Elegir foto"
            load()
        }
    }

    override fun onDestroyView() {
        super.onDestroyView()
        _binding = null
    }
}
