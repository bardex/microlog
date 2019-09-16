import Vue from 'vue';

import VueRouter from 'vue-router'
import InputsList from './inputs-list.vue'
import InputAdd from './input-add.vue'

Vue.use(VueRouter);

let router = new VueRouter({
    routes: [
        {path: '/inputs', component: InputsList},
        {path: '/input/add', component: InputAdd}
    ]
});


// 4. Создаём и монтируем корневой экземпляр приложения.
// Убедитесь, что передали экземпляр маршрутизатора в опции
// `router`, чтобы позволить приложению знать о его наличии.
new Vue({
    el:"#app",
    data: {
        message:"Hello"
    },
    router:router
});


