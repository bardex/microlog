<template>
    <form action="/input/add" method="post">


                <select name="protocol">
                    <option value="udp">UDP</option>
                    <option value="tcp">TCP</option>
                    <option value="http">HTTP</option>
                </select>


                <select   name="extractor">


                    <option value="JSON">JSON</option>

                    <option value="STRING">String</option>

                    <option value="ZLIB_JSON">zlib &#43; JSON (GrayLog)</option>


                </select>


                <input name="address" type="text" placeholder="Address (like :8081)">


                <button type="submit" class="btn btn-primary">Save</button>


    </form>
</template>

<script>
    export default {

    }
</script>

<style scoped>

</style>