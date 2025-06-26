import requests
import sys


def del_post(PORT, id):

    post = {"id": id}
    print(post)

    make_post = requests.delete(
        f"http://localhost:{PORT}/posts",
        # f"http://192.168.0.100:{PORT}/posts",
        json=post,
        headers={"Content-type": "application/json"},
    )

    print(f"Status code: {make_post.status_code}")
    print(f"Response: {make_post.text}")


if __name__ == "__main__":
    PORT = 8081
    id = sys.argv[1]  # 0 es el propio script
    # id = 7
    for _ in range(1):
        del_post(PORT, id)
