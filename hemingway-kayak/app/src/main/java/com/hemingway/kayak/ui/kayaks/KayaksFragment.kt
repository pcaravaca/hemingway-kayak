package com.hemingway.kayak.ui.kayaks

import android.os.Bundle
import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup
import android.widget.ArrayAdapter
import android.widget.TextView
import androidx.core.view.isVisible
import androidx.fragment.app.Fragment
import androidx.lifecycle.lifecycleScope
import com.google.android.material.dialog.MaterialAlertDialogBuilder
import com.hemingway.kayak.R
import com.hemingway.kayak.data.remote.KayakDto
import com.hemingway.kayak.data.remote.KayakRepository
import com.hemingway.kayak.data.remote.KayakRequest
import com.hemingway.kayak.databinding.FragmentKayaksBinding
import dagger.hilt.android.AndroidEntryPoint
import kotlinx.coroutines.launch
import javax.inject.Inject

/** Kayaks del usuario: listado, alta, edicion y borrado. */
@AndroidEntryPoint
class KayaksFragment : Fragment() {

    private var _binding: FragmentKayaksBinding? = null
    private val binding get() = _binding!!

    @Inject
    lateinit var kayakRepository: KayakRepository

    /** Id en edicion; null = alta. */
    private var editingId: String? = null

    override fun onCreateView(
        inflater: LayoutInflater,
        container: ViewGroup?,
        savedInstanceState: Bundle?
    ): View {
        _binding = FragmentKayaksBinding.inflate(inflater, container, false)
        return binding.root
    }

    override fun onViewCreated(view: View, savedInstanceState: Bundle?) {
        super.onViewCreated(view, savedInstanceState)
        setupSpinners()
        binding.addKayakButton.setOnClickListener {
            editingId = null
            clearForm()
            binding.kayakForm.isVisible = true
        }
        binding.cancelKayakButton.setOnClickListener {
            binding.kayakForm.isVisible = false
            editingId = null
        }
        binding.saveKayakButton.setOnClickListener { save() }
        load()
    }

    override fun onResume() {
        super.onResume()
        load()
    }

    private fun setupSpinners() {
        val types = resources.getStringArray(R.array.kayak_types)
        val materials = resources.getStringArray(R.array.kayak_materials)
        val propulsions = resources.getStringArray(R.array.kayak_propulsions)
        binding.kayakType.adapter = ArrayAdapter(requireContext(), android.R.layout.simple_spinner_dropdown_item, types)
        binding.kayakMaterial.adapter = ArrayAdapter(requireContext(), android.R.layout.simple_spinner_dropdown_item, materials)
        binding.kayakPropulsion.adapter = ArrayAdapter(requireContext(), android.R.layout.simple_spinner_dropdown_item, propulsions)
    }

    private fun load() {
        viewLifecycleOwner.lifecycleScope.launch {
            render(kayakRepository.all())
        }
    }

    private fun render(kayaks: List<KayakDto>) {
        binding.kayaksContainer.removeAllViews()
        binding.emptyView.isVisible = kayaks.isEmpty()
        val inflater = LayoutInflater.from(requireContext())
        kayaks.forEach { kayak ->
            val row = inflater.inflate(R.layout.item_kayak, binding.kayaksContainer, false)
            row.findViewById<TextView>(R.id.kayakName).text = kayak.name
            row.findViewById<TextView>(R.id.kayakSpecs).text = specs(kayak)
            row.findViewById<View>(R.id.editKayakButton).setOnClickListener { edit(kayak) }
            row.findViewById<View>(R.id.deleteKayakButton).setOnClickListener { confirmDelete(kayak) }
            binding.kayaksContainer.addView(row)
        }
    }

    private fun specs(k: KayakDto): String {
        val parts = mutableListOf<String>()
        k.type?.let { parts.add(it) }
        k.length?.takeIf { it > 0 }?.let { parts.add("${it.toInt()} cm") }
        k.width?.takeIf { it > 0 }?.let { parts.add("${it.toInt()} cm ancho") }
        k.weight?.takeIf { it > 0 }?.let { parts.add("${it.toInt()} kg") }
        k.maxCapacity?.takeIf { it > 0 }?.let { parts.add("carga máx. ${it.toInt()} kg") }
        k.material?.let { parts.add(it) }
        k.propulsion?.let { parts.add(it) }
        return parts.joinToString(" · ")
    }

    private fun edit(kayak: KayakDto) {
        editingId = kayak.id
        binding.kayakName.setText(kayak.name)
        select(binding.kayakType, R.array.kayak_types, kayak.type)
        select(binding.kayakMaterial, R.array.kayak_materials, kayak.material)
        select(binding.kayakPropulsion, R.array.kayak_propulsions, kayak.propulsion)
        binding.kayakLength.setText(kayak.length?.let { trim(it) } ?: "")
        binding.kayakWidth.setText(kayak.width?.let { trim(it) } ?: "")
        binding.kayakWeight.setText(kayak.weight?.let { trim(it) } ?: "")
        binding.kayakMaxCapacity.setText(kayak.maxCapacity?.let { trim(it) } ?: "")
        binding.kayakImageUrl.setText(kayak.imageUrl.orEmpty())
        binding.errorView.isVisible = false
        binding.kayakForm.isVisible = true
    }

    private fun select(spinner: android.widget.Spinner, arrayRes: Int, value: String?) {
        val options = resources.getStringArray(arrayRes)
        val index = options.indexOf(value)
        if (index >= 0) spinner.setSelection(index)
    }

    private fun trim(value: Double): String =
        if (value == value.toLong().toDouble()) value.toLong().toString() else value.toString()

    private fun clearForm() {
        binding.kayakName.text?.clear()
        binding.kayakType.setSelection(0)
        binding.kayakMaterial.setSelection(0)
        binding.kayakPropulsion.setSelection(0)
        binding.kayakLength.text?.clear()
        binding.kayakWidth.text?.clear()
        binding.kayakWeight.text?.clear()
        binding.kayakMaxCapacity.text?.clear()
        binding.kayakImageUrl.text?.clear()
        binding.errorView.isVisible = false
    }

    private fun save() {
        val name = binding.kayakName.text?.toString()?.trim().orEmpty()
        if (name.isEmpty()) {
            binding.errorView.text = "Ponele un nombre al kayak"
            binding.errorView.isVisible = true
            return
        }
        val request = KayakRequest(
            name = name,
            type = binding.kayakType.selectedItem?.toString(),
            length = binding.kayakLength.text?.toString()?.toDoubleOrNull(),
            width = binding.kayakWidth.text?.toString()?.toDoubleOrNull(),
            weight = binding.kayakWeight.text?.toString()?.toDoubleOrNull(),
            maxCapacity = binding.kayakMaxCapacity.text?.toString()?.toDoubleOrNull(),
            material = binding.kayakMaterial.selectedItem?.toString(),
            propulsion = binding.kayakPropulsion.selectedItem?.toString(),
            imageUrl = binding.kayakImageUrl.text?.toString()?.trim()?.ifEmpty { null }
        )
        binding.errorView.isVisible = false
        binding.saveKayakButton.isEnabled = false
        viewLifecycleOwner.lifecycleScope.launch {
            val id = editingId
            val ok = if (id == null) kayakRepository.create(request)
                     else kayakRepository.update(id, request)
            binding.saveKayakButton.isEnabled = true
            if (!ok) {
                binding.errorView.text = "No se pudo guardar el kayak"
                binding.errorView.isVisible = true
                return@launch
            }
            binding.kayakForm.isVisible = false
            editingId = null
            clearForm()
            load()
        }
    }

    private fun confirmDelete(kayak: KayakDto) {
        MaterialAlertDialogBuilder(requireContext())
            .setTitle("Eliminar kayak")
            .setMessage("¿Seguro querés eliminar ${kayak.name}? Esta acción no se puede deshacer.")
            .setNegativeButton("Cancelar", null)
            .setPositiveButton("Eliminar") { _, _ ->
                viewLifecycleOwner.lifecycleScope.launch {
                    if (!kayakRepository.delete(kayak.id)) {
                        binding.errorView.text = "No se pudo eliminar el kayak"
                        binding.errorView.isVisible = true
                    }
                    load()
                }
            }
            .show()
    }

    override fun onDestroyView() {
        super.onDestroyView()
        _binding = null
    }
}
