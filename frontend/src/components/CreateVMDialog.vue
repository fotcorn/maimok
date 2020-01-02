<template>
  <div>
    <v-dialog max-width="600px" v-model="open">
      <template v-slot:activator="{ on }">
        <v-btn bottom right color="primary" dark fab fixed v-on="on">
          <v-icon>fa-plus</v-icon>
        </v-btn>
      </template>
      <v-card>
        <v-card-title>
          <span class="headline">Create Virtual Machine</span>
        </v-card-title>
        <v-card-text>
          <v-form ref="form">
            <v-text-field
              label="Hostname"
              required
              v-model="hostname"
            ></v-text-field>
            <v-text-field
              label="IP Address"
              v-model="ipAddress"
              :prefix="ipAddressPrefix"
            ></v-text-field>
            <v-text-field label="Name" disabled :value="name"></v-text-field>
            <v-text-field
              label="RAM"
              required
              v-model="ram"
              suffix="MB"
            ></v-text-field>
            <v-text-field
              label="Disk space"
              required
              v-model="diskSpace"
              suffix="GB"
            ></v-text-field>
          </v-form>
        </v-card-text>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn color="primary" text @click="open = false">Cancel</v-btn>
          <v-btn color="primary" text @click="create">Create</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
    <v-snackbar v-model="statusOK" color="success">
      Successfully created virtual machine.
      <v-btn dark text @click="statusOK = false">
        Close
      </v-btn>
    </v-snackbar>
    <v-snackbar v-model="statusError" color="error">
      Failed to created virtual machine: {{ error }}
      <v-btn dark text @click="statusError = false">
        Close
      </v-btn>
    </v-snackbar>
  </div>
</template>

<script lang="ts">
import { Component, Prop, Vue } from "vue-property-decorator";

@Component
export default class CreateVMDialog extends Vue {
  open = false;
  ipAddress = "";
  hostname = "";
  ram = "";
  diskSpace = "";

  statusOK = false;
  statusError = false;
  error = "";

  readonly ipAddressPrefix = "192.168.0.";

  async create() {
    if (!(this.$refs.form as any).validate()) {
      return;
    }
    const response = await fetch("http://localhost:7000/vms", {
      method: "POST",
      headers: {
        "Content-Type": "application/json"
      },
      body: JSON.stringify({
        disk_space_gb: parseInt(this.diskSpace),
        ram_mb: parseInt(this.ram),
        name: this.name,
        hostname: this.hostname,
        ip_address: `${this.ipAddressPrefix}${this.ipAddress}`
      })
    });

    const responseJSON: any = await response.json();

    if (responseJSON.status == "ok") {
      this.statusOK = true;
      this.open = false;
    } else {
      this.error = responseJSON.error;
      this.statusError = true;
    }
  }

  get name() {
    if (this.ipAddress && this.hostname) {
      return `${this.ipAddress}_${this.hostname}`;
    } else {
      return "";
    }
  }
}
</script>
