<template>
  <v-card>
    <v-card-title>
      Virtual Machines
      <v-spacer></v-spacer>
      <v-text-field
        v-model="search"
        label="Search"
        single-line
        hide-details
      ></v-text-field>
    </v-card-title>
    <v-data-table
      :headers="headers"
      :items="vms"
      :search="search"
      disable-pagination
      hide-default-footer
    >
      <template v-slot:item.name="{ item }">
        <span class="fa-stack">
          <i class="fas fa-desktop fa-stack-2x"></i>
          <i
            v-if="item.running"
            class="fas fa-play fa-stack-1x stacked-icon success--text"
          ></i>
          <i
            v-else
            class="fas fa-pause fa-stack-1x stacked-icon grey--text text--darken-1"
          ></i>
        </span>

        <span class="ml-2">{{ item.name }}</span>
      </template>
      <template v-slot:item.memory="{ item }">
        {{ formatMemory(item.memory) }} GB
      </template>
    </v-data-table>
    <create-vm-dialog @created="load" />
  </v-card>
</template>

<script lang="ts">
import { Component, Prop, Vue } from "vue-property-decorator";
import CreateVMDialog from "../components/CreateVMDialog.vue";

@Component({
  components: {
    "create-vm-dialog": CreateVMDialog
  }
})
export default class VMList extends Vue {
  readonly headers = [
    { text: "Name", value: "name" },
    { text: "Memory", value: "memory" }
  ];

  vms = [];
  search = "";
  interval = null;

  async mounted() {
    await this.load();
    window.setInterval(this.load, 3000);
  }

  async load() {
    const response = await fetch("/api/vms");
    this.vms = await response.json();
  }

  formatMemory(memory: number) {
    return +(memory / 1024 / 1024).toFixed(2);
  }
}
</script>

<style lang="sass" scoped>
.stacked-icon
  font-size: 0.8em
  line-height: 2.0em
</style>
