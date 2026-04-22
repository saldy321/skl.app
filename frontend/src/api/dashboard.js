import api from './index';

export const getStats = (slug) => {
    return api.get(`/${slug}/dashboard`);
}