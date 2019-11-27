<template>
  <div class="header">
    <div class="control">
      <input type=file style="display: none;" ref="upload-file">
      <Button type="primary" icon="md-cloud-upload" @click="upload">上传</Button>
      <Button icon="md-add" style="margin-left: 10px;">新建文件夹</Button>
    </div>
    <div class="user">
      <Dropdown trigger="click" placement="bottom-end">
        <a href="javascript:void(0)">
            {{ user.username }}
            <Icon type="ios-arrow-down"></Icon>
        </a>
        <DropdownMenu slot="list">
            <DropdownItem>退出</DropdownItem>
        </DropdownMenu>
    </Dropdown>
    </div>
  </div>
</template>

<script>
import client from '../ajax-client'
import { mapState } from 'vuex'

export default {
  computed: {
    ...mapState([
      'user'
    ])
  },
  methods: {
    upload() {
      this.$refs['upload-file'].click()
    }
  },
  mounted() {
    this.$refs['upload-file'].onchange = e => {
      client.upload(e.target.files[0])
    }
  }
}
</script>

<style scoped>
.header {
  height: 100%;
  border-bottom: 1px #ddd solid;
  padding: 0 20px;
}
.control {
  display: inline-block;
}
.user {
  display: inline-block;
  position: absolute;
  right: 20px;
}
</style>