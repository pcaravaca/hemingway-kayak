package com.hemingway.kayak.ui.auth

import android.os.Bundle
import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup
import androidx.core.view.isVisible
import androidx.fragment.app.Fragment
import androidx.fragment.app.viewModels
import androidx.lifecycle.Lifecycle
import androidx.lifecycle.lifecycleScope
import androidx.lifecycle.repeatOnLifecycle
import androidx.navigation.fragment.findNavController
import com.hemingway.kayak.R
import com.hemingway.kayak.databinding.FragmentRegisterBinding
import dagger.hilt.android.AndroidEntryPoint
import kotlinx.coroutines.launch

@AndroidEntryPoint
class RegisterFragment : Fragment() {

    private var _binding: FragmentRegisterBinding? = null
    private val binding get() = _binding!!

    private val viewModel: RegisterViewModel by viewModels()

    override fun onCreateView(
        inflater: LayoutInflater,
        container: ViewGroup?,
        savedInstanceState: Bundle?
    ): View {
        _binding = FragmentRegisterBinding.inflate(inflater, container, false)
        return binding.root
    }

    override fun onViewCreated(view: View, savedInstanceState: Bundle?) {
        super.onViewCreated(view, savedInstanceState)
        binding.registerButton.setOnClickListener {
            viewModel.register(
                binding.nameInput.text?.toString().orEmpty(),
                binding.emailInput.text?.toString().orEmpty(),
                binding.passwordInput.text?.toString().orEmpty()
            )
        }
        binding.registerHaveAccount.setOnClickListener {
            findNavController().navigateUp()
        }
        observeState()
    }

    private fun observeState() {
        viewLifecycleOwner.lifecycleScope.launch {
            viewLifecycleOwner.repeatOnLifecycle(Lifecycle.State.STARTED) {
                viewModel.state.collect { state ->
                    binding.registerButton.isEnabled = !state.loading
                    binding.registerProgress.isVisible = state.loading
                    binding.registerError.isVisible = state.error != null
                    binding.registerError.text = state.error

                    if (state.success) {
                        goHome()
                    }
                }
            }
        }
    }

    private fun goHome() {
        // Registro y login quedan fuera de la pila: volver atrás no tiene que devolver acá.
        findNavController().navigate(
            R.id.homeFragment,
            null,
            androidx.navigation.navOptions {
                popUpTo(R.id.loginFragment) { inclusive = true }
            }
        )
    }

    override fun onDestroyView() {
        super.onDestroyView()
        _binding = null
    }
}
