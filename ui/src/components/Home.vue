<template>
  <Layout :style="{minHeight: '100vh'}">
    <Sider :style="{
      backgroundColor: '#fff',
      borderRight: '1px #dcdee2 solid'
    }">
      <div class="logo">Banana</div>
      <Menu active-name="1-1" width="auto">
        <MenuItem name="1-1">
          <Icon type="ios-navigate"></Icon>
          <span>全部</span>
        </MenuItem>
      </Menu>
    </Sider>
    <Layout>
      <Header class="header" :style="{backgroundColor: '#fff'}">
        Header
      </Header>
      <Content :style="{padding: '0 16px 16px'}">
        <Breadcrumb :style="{margin: '16px 0'}" separator=">">
          <BreadcrumbItem>全部</BreadcrumbItem>
          <BreadcrumbItem>folder</BreadcrumbItem>
        </Breadcrumb>
        <Card>
          <div style="height: auto">Content</div>
        </Card>
      </Content>
    </Layout>
  </Layout>
</template>

<script>
// import Sidebar from '@/components/Sidebar.vue'

export default {
  components: {
    //Sidebar
  },
  data() {
    return {

    }
  },
  computed: {
    files() {
      return this.$store.state.files
    }
  },
  methods: {
    refreshFiles() {
      this.$store.dispatch("updateFiles")
    }
  },
  created() {
    this.refreshFiles()
  }
}
</script>

<style .scoped>
.logo {
  margin: 11px 0;
  font-size: 1.75rem;
  text-align: center;
}
.header {
  border-bottom: 1px #ddd solid;
}
</style>

<style>
.ivu-menu-vertical.ivu-menu-light:after {
  content: none;
}

</style>