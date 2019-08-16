<template>
    <div v-if="inputs.length > 0">
        <table class="table table-sm table-bordered">
            <tr v-for="input in inputs">
                <td>#{{ input.id }}</td>
                <td>{{ input.protocol }} {{ input.addr }}</td>
                <td>
                    <div v-if="input.active">
                        <span class="badge badge-pill badge-success">ON</span>
                    </div>
                    <div v-else>
                        <span class="badge badge-pill badge-secondary">OFF</span>
                    </div>
                    <div class="text-danger" v-if="input.error.length > 0">
                        {{ input.error }}
                    </div>
                </td>
                <td>
                    <template v-if="input.active">
                        <button type="button" class="btn btn-danger" v-on:click="stop(input.stop_url)">Stop</button>
                    </template>
                    <template v-else>
                        <button type="button" class="btn btn-success" v-on:click="start(input.start_url)">Start</button>
                    </template>
                </td>
            </tr>
        </table>
    </div>
    <div v-else>
        No inputs
    </div>
</template>

<script>
    import axios from 'axios'

    export default {
        data: function () {
            return {
                inputs: []
            }
        },
        mounted: function () {
            this.reload()
        },
        methods: {
            "reload": function() {
                axios.get('/inputs')
                    .then(response => {
                        this.inputs = response.data
                    })
            },
            "stop": function (stopUrl) {
                if (confirm("Остановить этот слушатель ?")) {
                    axios.post(stopUrl).then(this.reload)
                }
            },
            "start": function (startUrl) {
                if (confirm("Запустить этот слушатель ?")) {
                    axios.post(startUrl).then(this.reload)
                }
            },
        }
    }
</script>

