<script setup>
import { defineRouteMeta, useRouter } from '@fesjs/fes'

defineRouteMeta({
  name: 'index',
  title: '首页',
})

const router = useRouter()

router.replace('/account')
</script>

<template>
  <h1>
    首页
  </h1>
</template>
