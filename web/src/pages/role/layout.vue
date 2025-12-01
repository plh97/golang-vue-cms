<script setup>
import { Page } from '@fesjs/fes'
import { FLayout, FMain } from '@fesjs/fes-design'
</script>

<template>
  <FLayout embedded>
    <FMain class="demo-main-wrapper">
      <div class="demo-main-page">
        <Page />
      </div>
    </FMain>
  </FLayout>
</template>

<style scoped>
.demo-main-wrapper {
  padding: 33px 24px;
}
.demo-main-page {
  padding: 15px 12px;
  box-sizing: border-box;
  background: #fff;
  border-radius: 6px;
  box-shadow: 0 2px 3px 0 rgba(0, 0, 0, 0.10);
  display: flex;
  flex-direction: column;
  overflow: auto;
  height: calc(100vh - 120px);
}
</style>
