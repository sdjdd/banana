import axios from 'axios'
import * as errors from './errors'
import store from './store'

const client = {
    axins: axios.create({
        auth: {
            username: '',
            password: ''
        }
    }),

    async login(username, password) {
        try {
            let config = { auth: { username, password } }
            await this.axins.get('/api/auth/verify', config)
            this.axins = axios.create(config)
        } catch (err) {
            if (err.response.status === 403) {
                throw errors.NOT_ALLOW
            }
            throw errors.INTERNAL
        }
    },

    async info() {
        try {
            let resp = await this.axins.get('/api/info')
            store.commit('setInfo', resp.data)
        } catch (err) {
            if (err.response.status === 403) {
                throw errors.NOT_ALLOW
            }
            throw errors.INTERNAL
        }
    },

    async files(path = '/') {
        try {
            let resp = await this.axins.get('/fs' + path)
            store.commit('setFiles', resp.data)
        } catch (err) {
            if (err.response.status === 403) {
                throw errors.NOT_ALLOW
            }
            throw errors.INTERNAL
        }
    }
}

export default client
