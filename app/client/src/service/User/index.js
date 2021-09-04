class UserService {
    getId() {
        return localStorage.getItem("user_id") ?? null;
    }

    setId(id) {
        localStorage.setItem("user_id", id);
    }
}

export default new UserService()