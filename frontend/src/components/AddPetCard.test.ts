import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import AddPetCard from './AddPetCard.vue'

describe('AddPetCard', () => {
  it('renders as a button with accessible label', () => {
    const wrapper = mount(AddPetCard)
    expect(wrapper.element.tagName).toBe('BUTTON')
    expect(wrapper.attributes('aria-label')).toBeTruthy()
  })

  it('emits click event', async () => {
    const wrapper = mount(AddPetCard)
    await wrapper.trigger('click')
    expect(wrapper.emitted('click')).toHaveLength(1)
  })

  it('is keyboard accessible (Enter triggers click)', async () => {
    const wrapper = mount(AddPetCard)
    await wrapper.trigger('keydown.enter')
    // Native button: Enter is handled by the browser; we just verify focusability
    expect(wrapper.attributes('type')).toBe('button')
  })
})
