import DynamicA from 'src/pages/DynamicA.vue'

const routes = [
  {
    path: '/',
    component: () => import('layouts/MainLayout.vue'),
    children: [
      {
        path: '',
        component: () => import('src/pages/HomeM.vue'),

        name: 'HomeM',
      },
      {
        path: '/dynamic',
        name: 'Dynamic',
        component: DynamicA,
      },
      {
        path: 'HotVedio',
        component: () => import('src/pages/HotVedio.vue'),
        name: 'HotVedio',
      },
      {
        path: '/videos/:id',
        component: () => import('src/pages/VedioCom.vue'),
        name: 'Video',
      },
      {
        path: '/personalHome/:userId',
        component: () => import('src/pages/PersonalHomepage.vue'),
        name: 'Personal',
      },
    ],
  },

  {
    path: '/:catchAll(.*)*',
    component: () => import('pages/ErrorNotFound.vue'),
  },
]

export default routes
