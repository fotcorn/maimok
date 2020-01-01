<template>
  <v-data-table :headers="headers" :items="vms" :items-per-page="100" />
</template>

<script lang="ts">
import { Component, Prop, Vue } from "vue-property-decorator";

@Component
export default class VMList extends Vue {
  readonly headers = [
    { text: "Running", value: "running" },
    { text: "Name", value: "name" },
    { text: "Memory", value: "memory" }
  ];

  vms = [];
  async mounted() {
    const response = await fetch("http://localhost:7000/vms");
    this.vms = await response.json();
  }
}
</script>
