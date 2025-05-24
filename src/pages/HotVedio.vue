<template>
  <q-layout>
    <q-page>
      <QToolbar>
        <div class="q-pa-md">
          <q-btn-dropdown color="pink" label="Dropdown Button" class="Sidebar">
            <q-list>
              <q-item clickable v-close-popup>
                <q-item-section>
                  <q-item-label>Photos</q-item-label>
                </q-item-section>
              </q-item>

              <q-item clickable v-close-popup>
                <q-item-section>
                  <q-item-label>Videos</q-item-label>
                </q-item-section>
              </q-item>

              <q-item clickable v-close-popup>
                <q-item-section>
                  <q-item-label>Articles</q-item-label>
                </q-item-section>
              </q-item>
            </q-list>
          </q-btn-dropdown>
        </div>
      </QToolbar>
    </q-page>
  </q-layout>
</template>

<style>
.Sidebar {
  position: fixed;
  top: 50%;
  left: 90%;
}
</style>
