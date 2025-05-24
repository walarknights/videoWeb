<template>
  <q-layout view="lHh lpr lFf" style="border: none" class="full-height">
    <q-header>
      <HeadeLead />
    </q-header>

    <q-page-container>
      <router-view></router-view>
    </q-page-container>
  </q-layout>
</template>

<script setup>
import HeadeLead from 'src/components/Heade-lead.vue'
</script>
