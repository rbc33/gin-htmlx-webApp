import requests
import os

path = "/Users/ric/code/rag/langchainpy/gocms/"
files = os.listdir(path)
for f in files:
    print(f)
with open(path + "gocms-cms-bugs.res.md", "r") as file:
    file_content = file.read()

print(file_content)
post = {
    "title": "Tend fire such.",
    "excerpt": "Your future respond stuff even behind probably. Bit agent community house. Of do film cold light standard.",
    "content": """{{img:4babde4e-f6af-489a-af86-04adedabd6f7.jpg}}---\n\nCare themselves course thank decide step I.\n\n
    ```commercial\nName fight the fly. Building already safe score experience five.\n```\n\n
    ---\n\n***\n\n```arm\nSenior conference sister author notice business.\n```\n\n
    1. Spring forward tonight pretty she phone.\n1. Game fast your debate activity amount.\n
    1. Particularly window thus toward which matter describe.\nProvide future hope plant statement.\n\n
    \n|More forward rise.|Way fast six decade remain once animal computer.|Teacher family interest win especially wish suggest.|\n
    |------------------|------------------------------------------------|----------------------------------------------------|\n
    |Door become indicate those.|Money wide suffer arm.|Pretty despite describe.|\n
    |Deep friend act tonight success east.|According than consumer.|Article large number rich issue interesting choice.|\n
    |Network area nor difficult for.|Safe education yeah almost training degree fall.|Line tough fast store million.|\n
    |While design blood audience film.|Interest positive around religious deep artist.|Light education seat oil.|\n\n
    \nLive large color hour official.\n-------------------------------\n\n```unit\nDiscover word way building. Recent some eye onto only.\n```\n\n
    **threat**\n\n |Attention debate almost draw.|Ten follow never eye me.|Determine cause though lot current.|\n
    |-----------------------------|------------------------|-----------------------------------|\n
    |Matter strategy entire choose grow agent become that.|Ok international art.|Especially dark actually economic from either young.|\n|Fact behind quality road up explain.|Knowledge foreign central future throughout hear get agent.|Place assume wait early short simple.|\n|Out glass already whatever prepare.|Hear magazine television.|For military Congress government should.|\n|Fast option and.|Better hour reach happy they listen game.|Officer result group degree citizen.|\n\n\n![Edge family with American her.](https://picsum.photos/426 "Age teacher high card speak section. Soldier because information employee season study.\nHappen start parent city easy beat. Eight reach exist gas.")\n\n[Successful cover produce certainly.](http://jones-arias.net/)\n\nPositive authority report month build.\n--------------------------------------\n\n**practice**\nkitchen\n> Performance media no station quickly scientist.\n\n
    [Oil plant happen will relationship.](https://smith.com/)\n\n""",
}

make_request = requests.post("http://localhost:8081/posts", json=post)
print(make_request.status_code, make_request.content)
