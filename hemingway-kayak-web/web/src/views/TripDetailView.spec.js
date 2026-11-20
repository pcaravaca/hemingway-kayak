import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'

// Leaflet no arranca en jsdom: se reemplaza por una versión que dibuja nada.
vi.mock('leaflet', () => {
  const chain = () => ({ addTo: () => chain(), bindPopup: () => chain(), remove: () => {} })
  return {
    default: {
      map: () => ({ setView: () => {}, remove: () => {}, fitBounds: () => {} }),
      tileLayer: () => ({ addTo: () => {} }),
      polyline: () => ({ addTo: () => {} }),
      latLngBounds: () => ({ pad: () => {} }),
      circleMarker: () => ({ addTo: () => chain(), bindPopup: () => chain() })
    }
  }
})

const { put, post, del } = vi.hoisted(() => ({
  put: vi.fn(),
  post: vi.fn(),
  del: vi.fn()
}))

vi.mock('axios', () => ({
  default: {
    get: vi.fn((url) => {
      if (url.endsWith('/points')) return Promise.resolve({ data: [] })
      return Promise.resolve({
        data: {
          id: 't1', title: 'Salida de prueba', description: 'desc', type: 'FISHING',
          status: 'COMPLETED', distance: 3.2, startTime: '2026-03-15T06:00:00Z',
          endTime: null, startLat: 9.9, startLon: -84.8
        }
      })
    }),
    put, post, delete: del
  }
}))

vi.mock('vue-router', () => ({
  useRoute: () => ({ params: { id: 't1' } }),
  useRouter: () => ({ replace: vi.fn() })
}))

import TripDetailView from '@/views/TripDetailView.vue'

const mountView = () =>
  mount(TripDetailView, {
    global: { stubs: { 'router-link': { template: '<a><slot /></a>' } } }
  })

const tick = () => new Promise((resolve) => setTimeout(resolve, 0))
const boton = (wrapper, texto) =>
  wrapper.findAll('button').find((b) => b.text() === texto)

describe('TripDetailView', () => {
  beforeEach(() => {
    put.mockReset().mockResolvedValue({})
    post.mockReset().mockResolvedValue({})
    del.mockReset().mockResolvedValue({})
    window.confirm = vi.fn(() => true)
  })

  it('muestra la salida y las acciones de editar, agregar punto y eliminar', async () => {
    const wrapper = mountView()
    await tick()
    expect(wrapper.text()).toContain('Salida de prueba')
    expect(wrapper.text()).toContain('Editar')
    expect(wrapper.text()).toContain('Eliminar')
    expect(wrapper.text()).toContain('Agregar punto')
  })

  it('guarda los cambios con PUT', async () => {
    const wrapper = mountView()
    await tick()
    await boton(wrapper, 'Editar').trigger('click')
    await boton(wrapper, 'Guardar').trigger('click')
    expect(put).toHaveBeenCalledWith(
      '/api/v1/trips/t1',
      expect.objectContaining({ title: 'Salida de prueba' })
    )
  })

  it('agrega un punto con POST', async () => {
    const wrapper = mountView()
    await tick()
    const inputs = wrapper.findAll('input[type=number]')
    await inputs[0].setValue('9.95')
    await inputs[1].setValue('-84.85')
    await boton(wrapper, 'Agregar punto').trigger('click')
    expect(post).toHaveBeenCalledWith(
      '/api/v1/trips/t1/points',
      expect.objectContaining({ latitude: 9.95, longitude: -84.85 })
    )
  })

  it('borra la salida con DELETE', async () => {
    const wrapper = mountView()
    await tick()
    await boton(wrapper, 'Eliminar').trigger('click')
    expect(del).toHaveBeenCalledWith('/api/v1/trips/t1')
  })
})
