<template>
  <div>
    <template v-if="running === false">
      <v-btn icon text @click="action('start')">
        <v-icon small>fa-play</v-icon>
      </v-btn>
    </template>
    <template v-else>
      <v-btn icon text @click="action('stop')">
        <v-icon small>fa-power-off</v-icon>
      </v-btn>

      <v-btn icon text @click="action('restart')">
        <v-icon small>fa-sync</v-icon>
      </v-btn>

      <v-menu offset-y>
        <template v-slot:activator="{ on }">
          <v-btn text icon color="primary" dark v-on="on">
            ...
          </v-btn>
        </template>
        <v-list>
          <v-list-item @click="action('forcestop')">
            <v-list-item-icon>
              <v-icon small>fa-exclamation</v-icon>
              <v-icon small>fa-power-off</v-icon>
            </v-list-item-icon>
            <v-list-item-content>
              Force shutdown
            </v-list-item-content>
          </v-list-item>
          <v-list-item @click="action('forcerestart')">
            <v-list-item-icon>
              <v-icon small>fa-exclamation</v-icon>
              <v-icon small>fa-sync</v-icon>
            </v-list-item-icon>
            <v-list-item-content>
              Force restart
            </v-list-item-content>
          </v-list-item>
        </v-list>
      </v-menu>
    </template>
    <v-snackbar v-model="statusOK" color="success">
      Successfully {{ messageOK }} virtual machine.
      <v-btn dark text @click="statusOK = false">
        Close
      </v-btn>
    </v-snackbar>
    <v-snackbar v-model="statusError" color="error">
      Failed to {{ messageError }} virtual machine: {{ error }}
      <v-btn dark text @click="statusError = false">
        Close
      </v-btn>
    </v-snackbar>
  </div>
</template>

<script lang="ts">
import { Component, Prop, Vue } from "vue-property-decorator";

@Component
export default class VMAction extends Vue {
  @Prop(Boolean) readonly running!: boolean;
  @Prop(String) readonly name!: boolean;

  statusOK = false;
  statusError = false;
  messageOK = "";
  messageError = "";
  error = "";

  async action(action: string) {
    const response = await fetch(`/api/vm/${action}`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json"
      },
      body: JSON.stringify({
        name: this.name
      })
    });

    const responseJSON: any = await response.json();

    if (responseJSON.status == "ok") {
      switch (action) {
        case "start":
          this.messageOK = "started";
          break;
        case "stop":
          this.messageOK = "stopped";
          break;
        case "forcestop":
          this.messageOK = "force stopped";
          break;
        case "restart":
          this.messageOK = "restarted";
          break;
        case "forcerestart":
          this.messageOK = "force restarted";
          break;
      }
      this.statusOK = true;
    } else {
      switch (action) {
        case "start":
          this.messageError = "start";
          break;
        case "stop":
          this.messageError = "stop";
          break;
        case "forcestop":
          this.messageError = "force stop";
          break;
        case "restart":
          this.messageError = "restart";
          break;
        case "forcerestart":
          this.messageError = "force";
          break;
      }
      this.error = responseJSON.error;
      this.statusError = true;
    }
  }
}
</script>
