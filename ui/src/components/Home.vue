<template>
  <Layout :style="{minHeight: '100vh'}">
    <Sider :style="{
      backgroundColor: '#fff',
      borderRight: '1px #dcdee2 solid'
    }">
      <MySider />
    </Sider>
    <Layout>
      <Header :style="{backgroundColor: '#fff', padding: 0}">
        <MyHeader />
      </Header>
      <Content :style="{padding: '0 16px 16px'}">
        <Breadcrumb :style="{margin: '16px 0'}" separator=">">
          <BreadcrumbItem>全部</BreadcrumbItem>
          <BreadcrumbItem>folder</BreadcrumbItem>
        </Breadcrumb>
        <Card dis-hover>
          <div style="height: auto">
            <Table :columns="fileColumns" :data="files">
              <template slot-scope="{ row }" slot="name">
                <a v-if="row.isDir" href="javascript:;" @click="move(row.name)">{{ row.name }}</a>
                <span v-else>{{ row.name }}</span>
              </template>
            </Table>
          </div>
        </Card>
      </Content>
    </Layout>
  </Layout>
</template>

<script>
import client from '../ajax-client'
import * as errors from '../errors'
import MyHeader from '@/components/Header.vue'
import MySider from '@/components/Sider.vue'

export default {
  components: {
    MyHeader,
    MySider
  },
  data() {
    return {
      fileColumns: [
        {
          title: '名称',
          slot: 'name'
        },
        {
          title: '上次修改时间',
          key: 'modTime'
        },
        {
          title: '大小',
          key: 'size'
        }
      ]
    }
  },
  computed: {
    files() {
      return this.$store.state.files.map(f => {
        if (f.isDir) {
          f.size = '-'
        }
        return f
      })
    },
    currentPath() {
      return this.$router.currentRoute.path
    }
  },
  methods: {
    refreshInfo() {
      client.info()
    },
    move(to) {
      this.$router.push('/'+to)
    }
  },
  async created() {
    try {
      await client.login('sdjdd', 'secret')
      this.refreshInfo()
      client.files(this.currentPath)
    } catch (err) {
      if (err === errors.NOT_ALLOW) {
        this.$router.push('/login')
      }
    }
  },
  beforeRouteUpdate: (to, from, next) => {
    client.files(to.fullPath)
    next()
  }
}
</script>

<style>
.logo {
  margin: 11px 0;
  font-size: 1.75rem;
  text-align: center;
}
</style>

<style>
.ivu-menu-vertical.ivu-menu-light:after {
  content: none;
}

</style>