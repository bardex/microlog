<template>
    <div>
        <div class="row">
            <div class="col s12 m6">
                <h4>Inputs <a href="/input/add" class="btn btn-floating  right"><i class="material-icons">add</i></a></h4>
            </div>
            
        </div>

        <div v-if="inputs.length > 0">
            <div class="row" v-for="input in inputs">
                <div class="col s12 m6">
                    <div class="card">
                        <div class="card-content">
                            <div class="card-title">

                                <span class="grey-text text-lighten-1"><small>#{{ input.id }}</small></span>

                                 {{ input.protocol }} {{ input.addr }}

                                <template v-if="input.active">
                                    <span class="badge green darken-1 white-text" title="This input is working now">ON</span>
                                </template>
                                <template v-else>
                                    <span class="badge grey white-text" title="This input is off now">OFF</span>
                                </template>
                            </div>
                            <div>
                                <div>
                                    Extractor: {{input.extractor}}
                                </div>
                                <div class="red-text text-accent-4" v-if="input.error.length > 0">
                                    Error: {{ input.error }}
                                </div>
                            </div>
                        </div>
                        <div class="card-action">
                            <template v-if="input.active">
                                <button type="button"
                                        class="btn red darken-1"
                                        v-on:click="stop(input.stop_url)"
                                        title="Click for stop this input">
                                    Stop
                                </button>
                            </template>
                            <template v-else>
                                <button type="button"
                                        class="btn red darken-1"
                                        v-on:click="del(input.delete_url)"
                                        title="Click for delete this input">Delete</button>

                                <button type="button"
                                        class="btn"
                                        v-on:click="start(input.start_url)"
                                        title="Click for start this input">Start</button>
                            </template>
                        </div>
                    </div>
                </div>
            </div>
        </div>
        <div v-else>
            No inputs
        </div>
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
            "reload": function () {
                axios.get('/api/inputs')
                    .then(response => {
                        this.inputs = response.data
                    })
            },
            "stop": function (stopUrl) {
                if (confirm("Stop it ?")) {
                    axios.post(stopUrl).then(this.reload)
                }
            },
            "start": function (startUrl) {
                if (confirm("Start it ?")) {
                    axios.post(startUrl).then(this.reload)
                }
            },
            "del": function (deleteUrl) {
                if (confirm("Delete this ?")) {
                    axios.post(deleteUrl).then(this.reload)
                }
            }
        }
    }
</script>

