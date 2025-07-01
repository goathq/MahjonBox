<template>
  <div class="control-panel">
    <button @click="toggleService('skin')">换肤{{ skinEnabled ? '关闭' : '开启' }}</button>
    <button @click="toggleService('helper')">牌效助手{{ helperEnabled ? '关闭' : '开启' }}</button>
  </div>
</template>

<script>
export default {
  methods: {
    async toggleService(service) {
      const result = await window.backend.ToggleService(service)
      this[`${service}Enabled`] = !this[`${service}Enabled`]
    }
  }
}
</script>