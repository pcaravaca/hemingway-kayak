package com.hemingway.kayak.ui.home

import android.os.Bundle
import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup
import androidx.fragment.app.Fragment
import androidx.navigation.fragment.findNavController
import com.hemingway.kayak.R
import com.hemingway.kayak.databinding.FragmentHomeBinding
import dagger.hilt.android.AndroidEntryPoint

@AndroidEntryPoint
class HomeFragment : Fragment() {

    private var _binding: FragmentHomeBinding? = null
    private val binding get() = _binding!!

    override fun onCreateView(
        inflater: LayoutInflater,
        container: ViewGroup?,
        savedInstanceState: Bundle?
    ): View {
        _binding = FragmentHomeBinding.inflate(inflater, container, false)
        return binding.root
    }

    override fun onViewCreated(view: View, savedInstanceState: Bundle?) {
        super.onViewCreated(view, savedInstanceState)
        setupClickListeners()
    }

    private fun setupClickListeners() {
        binding.weatherButton.setOnClickListener {
            findNavController().navigate(R.id.action_home_to_weather)
        }

        binding.fishingButton.setOnClickListener {
            findNavController().navigate(R.id.action_home_to_fishing)
        }

        binding.mapButton.setOnClickListener {
            findNavController().navigate(R.id.action_home_to_map)
        }

        binding.historyButton.setOnClickListener {
            findNavController().navigate(R.id.historyFragment)
        }

        binding.settingsButton.setOnClickListener {
            findNavController().navigate(R.id.action_home_to_settings)
        }
    }

    override fun onDestroyView() {
        super.onDestroyView()
        _binding = null
    }
}
