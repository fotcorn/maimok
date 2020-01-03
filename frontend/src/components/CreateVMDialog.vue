<template>
  <div>
    <v-btn bottom right color="primary" dark fab fixed @click="show">
      <v-icon>fa-plus</v-icon>
    </v-btn>
    <v-dialog max-width="600px" v-model="open">
      <v-card :loading="loading">
        <v-card-title>
          <span class="headline">Create Virtual Machine</span>
        </v-card-title>
        <v-card-text>
          <v-form ref="form">
            <v-text-field
              label="Hostname"
              v-model="hostname"
              :rules="rules.hostname"
            ></v-text-field>
            <v-text-field
              label="IP address"
              v-model="ipAddress"
              :prefix="ipAddressPrefix"
              :rules="rules.ipAddress"
            ></v-text-field>
            <v-text-field label="Name" disabled :value="name"></v-text-field>
            <v-text-field
              label="RAM"
              v-model="ram"
              suffix="MB"
              :rules="rules.ram"
            ></v-text-field>
            <v-text-field
              label="Disk space"
              v-model="diskSpace"
              suffix="GB"
              :rules="rules.diskSpace"
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
  // form values
  ipAddress = "";
  hostname = "";
  ram = "";
  diskSpace = "";

  readonly ipAddressPrefix = "192.168.0.";

  // form validation

  isInteger(num: string) {
    return /^(0|[1-9]\d*)$/.test(num);
  }

  readonly rules = {
    hostname: [(v: string) => v.length != 0 || "Hostname is required"],
    ipAddress: [
      (v: string) => v.length != 0 || "IP address is required",
      (v: string) => this.isInteger(v) || "IP address is not a number"
    ],
    ram: [
      (v: string) => v.length != 0 || "RAM is required",
      (v: string) => this.isInteger(v) || "IP address is not a positive number",
      (v: string) => parseInt(v) >= 128 || "RAM must be at least 128 MB"
    ],
    diskSpace: [
      (v: string) => v.length != 0 || "Disk space is required",
      (v: string) => this.isInteger(v) || "Disk space is not a positive number",
      (v: string) => parseInt(v) >= 5 || "Disk space must be at least 5 GB"
    ]
  };

  // ui state
  open = false;
  statusOK = false;
  statusError = false;
  error = "";
  loading = false;

  show() {
    if (this.$refs.form) {
      (this.$refs.form as any).resetValidation();
    }
    this.open = true;
  }

  async create() {
    if (!(this.$refs.form as any).validate()) {
      return;
    }
    this.loading = true;
    const response = await fetch("/api/vms", {
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
    this.loading = false;

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
