package com.hemingway.kayak.ui.profile

import android.os.Bundle
import android.view.Gravity
import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup
import android.widget.TextView
import androidx.core.view.isVisible
import androidx.fragment.app.Fragment
import androidx.lifecycle.lifecycleScope
import androidx.navigation.fragment.findNavController
import com.hemingway.kayak.R
import com.hemingway.kayak.data.remote.StatsDto
import com.hemingway.kayak.data.remote.UserDto
import com.hemingway.kayak.data.remote.UserRepository
import com.hemingway.kayak.data.remote.AuthRepository
import com.hemingway.kayak.databinding.FragmentProfileBinding
import dagger.hilt.android.AndroidEntryPoint
import kotlinx.coroutines.launch
import javax.inject.Inject

@AndroidEntryPoint
class ProfileFragment : Fragment() {

    private var _binding: FragmentProfileBinding? = null
    private val binding get() = _binding!!

    @Inject
    lateinit var userRepository: UserRepository

    @Inject
    lateinit var authRepository: AuthRepository

    override fun onCreateView(
        inflater: LayoutInflater,
        container: ViewGroup?,
        savedInstanceState: Bundle?
    ): View {
        _binding = FragmentProfileBinding.inflate(inflater, container, false)
        return binding.root
    }

    override fun onViewCreated(view: View, savedInstanceState: Bundle?) {
        super.onViewCreated(view, savedInstanceState)
        binding.logoutButton.setOnClickListener { logout() }
        load()
    }

    private fun load() {
        binding.progressBar.isVisible = true
        binding.errorView.isVisible = false
        viewLifecycleOwner.lifecycleScope.launch {
            val me = userRepository.me()
            val stats = userRepository.stats()
            binding.progressBar.isVisible = false
            if (me == null) {
                binding.errorView.text = "No se pudo cargar tu perfil"
                binding.errorView.isVisible = true
                return@launch
            }
            render(me, stats)
        }
    }

    private fun render(me: UserDto, stats: StatsDto?) {
        binding.nameText.text = me.displayName ?: me.username ?: "Kayakista"
        binding.usernameText.text = "@${me.username ?: ""}"
        binding.emailText.text = me.email
        binding.bioText.text = me.bio.orEmpty()

        val rows = listOf(
            "Salidas" to stats?.trips,
            "Sesiones en vivo" to stats?.sessions,
            "Spots" to stats?.spots,
            "Torneos" to stats?.tournaments
        )

        val grid = binding.statsGrid
        grid.removeAllViews()
        for ((label, value) in rows) {
            val cell = TextView(requireContext()).apply {
                text = "${value ?: 0}\n$label"
                textSize = 16f
                setTextColor(0xFFFFFFFF.toInt())
                gravity = Gravity.CENTER
                setPadding(24, 24, 24, 24)
                layoutParams = android.widget.GridLayout.LayoutParams().apply {
                    columnSpec = android.widget.GridLayout.spec(android.widget.GridLayout.UNDEFINED, 1f)
                }
            }
            grid.addView(cell)
        }
    }

    private fun logout() {
        authRepository.logout()
        findNavController().navigate(
            R.id.loginFragment,
            null,
            androidx.navigation.navOptions {
                popUpTo(R.id.nav_graph) { inclusive = true }
            }
        )
    }

    override fun onDestroyView() {
        super.onDestroyView()
        _binding = null
    }
}
