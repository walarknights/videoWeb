<template>
  <q-btn flat round dense icon="sym_o_menu">
    <q-menu>
      <q-list style="min-width: 100px">
        <q-item clickable v-close-popup>
          <q-item-section>关注的人</q-item-section>
        </q-item>
        <q-item clickable v-close-popup>
          <q-item-section>消息</q-item-section>
        </q-item>
        <q-item clickable v-close-popup>
          <q-item-section>创作中心</q-item-section>
        </q-item>
        <q-separator />
      </q-list>
    </q-menu>
  </q-btn>
</template>
<script setup></script>
