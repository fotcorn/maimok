import Vue from "vue";
import VueRouter from "vue-router";
import VMList from "../views/VMList.vue";

Vue.use(VueRouter);

const routes = [
  {
    path: "/",
    name: "vmlist",
    component: VMList
  }
];

const router = new VueRouter({
  mode: "history",
  base: process.env.BASE_URL,
  routes
});

export default router;
