<template>
  <div>
    <v-card class="mx-auto my-12" max-width="400">
      <form enctype="multipart/form-data" action="http://localhost:8081/upload" method="POST">
        <v-file-input show-size label="Upload Demo" name="myFile"/>
        <div class="text-center">
          <v-btn type="submit">Submit</v-btn>
        </div>
      </form>
    </v-card>
  </div>
</template>

<script>
export default {
  
}
</script>