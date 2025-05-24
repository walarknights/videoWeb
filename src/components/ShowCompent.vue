<template>
  <div>
    <q-card>
      <q-tabs
        v-model="tab"
        dense
        class="text-grey"
        active-color="primary"
        indicator-color="primary"
        align="justify"
      >
        <q-tab name="video" label="视频主页" />
        <q-tab name="focus" label="关注列表" />
        <q-tab name="fans" label="粉丝列表" />
      </q-tabs>
      <q-separator />
      <q-tab-panels v-model="tab" animated>
        <q-tab-panel name="video" class="q-pa-none" style="height: 500px">
          <q-splitter v-model="splitterModel">
            <template v-slot:before>
              <q-tabs
                v-model="innerTab"
                vertical
                class="custom-tabs column justify-evenly"
                active-color="negative"
                indicator-color="negative"
                style="height: 500px"
              >
                <q-tab name="innervideo" icon="sym_o_token" label="代表作" />
                <q-tab name="innerfocus" icon="sym_o_favorite" label="最近收藏" />
                <q-tab name="innerfans" icon="sym_o_recommend" label="最近点赞" />
              </q-tabs>
            </template>

            <template v-slot:after>
              <q-tab-panels
                v-model="innerTab"
                animated
                transition-prev="jump-up"
                transition-next="jump-up"
                style="height: 500px"
              >
                <q-tab-panel name="innervideo" style="height: 100%">
                  <div class="text-h4 q-mb-md">video</div>
                  <p>
                    Lorem ipsum dolor sit, amet consectetur adipisicing elit. Quis praesentium
                    cumque magnam odio iure quidem, quod illum numquam possimus obcaecati commodi
                    minima assumenda consectetur culpa fuga nulla ullam. In, libero.
                  </p>
                </q-tab-panel>

                <q-tab-panel name="innerfocus">
                  <div class="text-h4 q-mb-md">focus</div>
                  <p>
                    Lorem ipsum dolor sit, amet consectetur adipisicing elit. Quis praesentium
                    cumque magnam odio iure quidem, quod illum numquam possimus obcaecati commodi
                    minima assumenda consectetur culpa fuga nulla ullam. In, libero.
                  </p>
                </q-tab-panel>

                <q-tab-panel name="innerfans">
                  <div class="text-h4 q-mb-md">fans</div>
                  <p>
                    Lorem ipsum dolor sit, amet consectetur adipisicing elit. Quis praesentium
                    cumque magnam odio iure quidem, quod illum numquam possimus obcaecati commodi
                    minima assumenda consectetur culpa fuga nulla ullam. In, libero.
                  </p>
                </q-tab-panel>
              </q-tab-panels>
            </template>
          </q-splitter>
        </q-tab-panel>

        <q-tab-panel name="focus" style="height: 500px">
          <FocusList />
        </q-tab-panel>

        <q-tab-panel name="fans" style="height: 500px">
          <div class="text-h6">fans</div>
          Lorem ipsum dolor sit amet consectetur adipisicing elit.
        </q-tab-panel>
      </q-tab-panels>
    </q-card>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import FocusList from './FocusList.vue'
const tab = ref('video')
const innerTab = ref('innervideo')
const splitterModel = ref('20')
</script>
<style>
.vertical-tabs .q-tabs__content {
  display: flex;
  flex-direction: column;
  justify-content: space-between; /* 垂直方向均匀分布 */
}

.vertical-tabs .q-tab {
  flex: 1; /* 每个 tab 占据相等高度 */
  display: flex;
  flex-direction: column;
  align-items: center; /* 水平居中 */
  justify-content: center; /* 垂直居中 */
}

.custom-tabs .q-tab {
  color: rgb(14, 165, 224) !important; /* 未选中tab文字和图标为灰色 */
}
.custom-tabs .q-tab--active {
  color: #f44336 !important; /* 选中tab为红色，可根据active-color设置 */
}
</style>
