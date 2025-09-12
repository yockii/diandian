import { createApp } from 'vue'
import "@wailsio/runtime"
// Vuetify
import 'vuetify/styles'
import { createVuetify } from 'vuetify'
import * as components from 'vuetify/components'
import * as directives from 'vuetify/directives'
import { aliases, fa } from 'vuetify/iconsets/fa'
import App from './App.vue'

import '@fortawesome/fontawesome-free/css/all.css'
import 'unfonts.css'

const vuetify = createVuetify({
  components,
  directives,
  icons: {
    defaultSet: 'fa',
    aliases,
    sets: {
      fa,
    },
  }
})

createApp(App).use(vuetify).mount('#app')
