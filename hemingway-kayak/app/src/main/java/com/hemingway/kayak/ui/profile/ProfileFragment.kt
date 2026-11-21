package com.hemingway.kayak.ui.profile

import android.graphics.BitmapFactory
import android.net.Uri
import android.os.Bundle
import android.view.Gravity
import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup
import android.widget.TextView
import androidx.activity.result.contract.ActivityResultContracts
import androidx.core.view.isVisible
import androidx.fragment.app.Fragment
import androidx.lifecycle.lifecycleScope
import androidx.navigation.fragment.findNavController
import com.hemingway.kayak.BuildConfig
import com.hemingway.kayak.R
import com.hemingway.kayak.data.remote.PhotoRepository
import com.hemingway.kayak.data.remote.StatsDto
import com.hemingway.kayak.data.remote.UserDto
import com.hemingway.kayak.data.remote.UserRepository
import com.hemingway.kayak.data.remote.AuthRepository
import com.hemingway.kayak.databinding.FragmentProfileBinding
import dagger.hilt.android.AndroidEntryPoint
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.launch
import kotlinx.coroutines.withContext
import java.net.URL
import javax.inject.Inject

@AndroidEntryPoint
class ProfileFragment : Fragment() {

    private var _binding: FragmentProfileBinding? = null
    private val binding get() = _binding!!

    @Inject
    lateinit var userRepository: UserRepository

    @Inject
    lateinit var authRepository: AuthRepository

    @Inject
    lateinit var photoRepository: PhotoRepository

    private val pickAvatar = registerForActivityResult(
        ActivityResultContracts.GetContent()
    ) { uri ->
        if (uri != null) changeAvatar(uri)
    }

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
        binding.changeAvatarButton.setOnClickListener { pickAvatar.launch("image/*") }
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
        loadAvatar(me.avatarUrl)

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

    // La foto se sube sola al elegirla (PhotoRepository) y se guarda en el perfil.
    private fun changeAvatar(uri: Uri) {
        binding.changeAvatarButton.isEnabled = false
        binding.errorView.isVisible = false
        viewLifecycleOwner.lifecycleScope.launch {
            val me = userRepository.me()
            val url = photoRepository.upload(uri)
            if (url == null) {
                binding.errorView.text = "No se pudo subir la foto"
                binding.errorView.isVisible = true
            } else {
                val updated = userRepository.update(
                    id = me?.id ?: return@launch,
                    displayName = me.displayName,
                    bio = me.bio,
                    trackingEnabled = me.trackingEnabled ?: false,
                    avatarUrl = url
                )
                if (updated == null) {
                    binding.errorView.text = "No se pudo guardar la foto"
                    binding.errorView.isVisible = true
                } else {
                    loadAvatar(url)
                }
            }
            binding.changeAvatarButton.isEnabled = true
        }
    }

    // Sin librería de imágenes: se decodifica en IO y se asigna al ImageView.
    // El backend devuelve la URL relativa (/uploads/...), así que se resuelve
    // contra la base de la API; una URL absoluta se usa tal cual.
    private fun loadAvatar(url: String?) {
        if (url.isNullOrBlank()) {
            binding.avatarImage.setImageDrawable(null)
            return
        }
        val absolute = if (url.startsWith("http")) url
        else BuildConfig.API_BASE_URL.trimEnd('/') + "/" + url.trimStart('/')
        viewLifecycleOwner.lifecycleScope.launch {
            val bmp = withContext(Dispatchers.IO) {
                runCatching {
                    URL(absolute).openStream().use { BitmapFactory.decodeStream(it) }
                }.getOrNull()
            }
            if (bmp != null) binding.avatarImage.setImageBitmap(bmp)
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
