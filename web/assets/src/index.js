import Vue from 'vue'
import Inputs from './inputs-list.vue'



// 4. Создаём и монтируем корневой экземпляр приложения.
// Убедитесь, что передали экземпляр маршрутизатора в опции
// `router`, чтобы позволить приложению знать о его наличии.
const app = new Vue({
    render: h => h(Inputs)
}).$mount('#app');