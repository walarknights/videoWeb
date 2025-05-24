<template>
  <div class="q-pa-md" style="padding: 0%">
    <q-toolbar class="row no-wrap" style="justify-content: space-around; height: 100%; padding: 0%">
      <q-avatar @click="toHome">
        <img src="https://cdn.quasar.dev/logo-v2/svg/logo.svg" />
      </q-avatar>

      <q-tabs style="margin-left: 0px; margin-right: 75px !important">
        <q-tab name="images" label="热门" @click="toHot" />
        <q-tab name="videos" label="Videos" />
        <q-tab name="articles" label="Articles" />
      </q-tabs>
      <q-input label-color="primary" outlined v-model="text" style="flex-grow: 1">
        <q-icon name="sym_o_search" color="white" style="height: 100%" size="40px" />
      </q-input>
      <q-tabs v-model="tab">
        <q-tab name="images" label="Images" />
        <q-tab name="videos" label="Videos" />
        <q-tab name="articles" label="Articles" />
      </q-tabs>

      <div class="avatar-container">
        <AvatarLogin style="margin-right: 50%" />
      </div>
      <Menu_avatar />
    </q-toolbar>
  </div>
</template>

<script setup>
import Menu_avatar from './Menu_avatar.vue'
import AvatarLogin from './AvatarLogin.vue'
import { useRouter } from 'vue-router'
const router = useRouter()

const toHome = () => {
  router.push({ name: 'HomeM' })
}
const toHot = () => {
  router.push({ name: 'HotVedio' })
}
</script>
