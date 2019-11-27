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
      <Content :style="{
        padding: '0 16px 16px',
        display: 'flex',
        flexDirection: 'column'
      }">
        <Breadcrumb :style="{margin: '16px 0'}" separator=">">
          <BreadcrumbItem v-for="(dir, index) in dirStack" :key="dir.path">
            <a v-if="index < dirStack.length-1" href="javascript:;" @click="moveAbs(dir.path)">{{ dir.name }}</a>
            <span v-else>{{ dir.name }}</span>
          </BreadcrumbItem>
        </Breadcrumb>
        <Card dis-hover style="height: 100%;">
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
import { humanSize, humanDate } from '../utils'
import { mapState } from 'vuex'

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
          key: 'modTimeStr'
        },
        {
          title: '大小',
          key: 'size'
        }
      ],
    }
  },
  computed: {
    files() {
      return this.$store.state.files.map(f => {
        if (f.isDir) {
          f.size = '-'
        } else {
          f.size = humanSize(f.size)
        }
        f.modTime = new Date(f.modTime)
        f.modTimeStr = humanDate(f.modTime)
        return f
      }).sort((a, b) => {
        if (a.isDir && b.isDir || !a.isDir && !b.isDir) {
          return b.modTime - a.modTime
        } else if (a.isDir) {
          return -1
        } else if (b.isDir) {
          return 1
        }
      })
    },
    dirStack() {
      let stack = [ { name: '全部', path: '' } ]
      if (this.path !== '/') {
        this.path.split('/').slice(1).forEach(v => stack.push({ name: v }))
      }
      for (let i = 1; i < stack.length; ++i) {
        stack[i].path = stack[i-1].path + '/' + stack[i].name
      }
      stack[0].path = '/'
      return stack
    },
    ...mapState([
      'path'
    ])
  },
  methods: {
    refreshInfo() {
      client.info()
    },
    move(to) {
      this.$router.push((this.path === '/' ? '' : this.path + '/') + to)
    },
    moveAbs(to) {
      this.$router.push(to)
    }
  },
  async created() {
    try {
      await client.login('sdjdd', 'secret')
      this.refreshInfo()

      this.$store.commit('path', this.$router.currentRoute.path)
      client.files(this.path)
    } catch (err) {
      if (err === errors.NOT_ALLOW) {
        this.$router.push('/login')
      }
    }
  },
  beforeRouteUpdate(to, from, next) {
    this.$store.commit('path', to.fullPath)
    client.files(to.fullPath).then(next)
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