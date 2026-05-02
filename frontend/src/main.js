import { createApp } from 'vue'
import { createPinia } from 'pinia'
import App from './App.vue'
import router from './router'

import { 
  Button, NavBar, Tabbar, TabbarItem, Swipe, SwipeItem, 
  List, PullRefresh, Search, Empty,
  Form, Field, CellGroup,
  Image as VanImage, Icon, Popup,
  Loading, Cell, Grid, GridItem,
  Tag, Divider, Tabs, Tab,
  Stepper
} from 'vant'

import 'vant/lib/index.css'
import './styles/index.less'

const app = createApp(App)

app.use(createPinia())
app.use(router)

app.use(Button)
   .use(NavBar)
   .use(Tabbar)
   .use(TabbarItem)
   .use(Swipe)
   .use(SwipeItem)
   .use(List)
   .use(PullRefresh)
   .use(Search)
   .use(Empty)
   .use(Form)
   .use(Field)
   .use(CellGroup)
   .use(VanImage)
   .use(Icon)
   .use(Popup)
   .use(Loading)
   .use(Cell)
   .use(Grid)
   .use(GridItem)
   .use(Tag)
   .use(Divider)
   .use(Tabs)
   .use(Tab)
   .use(Stepper)

app.mount('#app')
