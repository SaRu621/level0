import tkinter as tk
import requests
import json

def send_request():
    url = 'http://localhost:9000/order/' + entry_url.get()

    response = requests.get(url)  # выполняем GET-запрос

    response_data = response.json()  # получаем данные из ответа в формате JSON

    entry_response.delete("1.0", tk.END)  # очищаем поле ответа
    entry_response.insert(tk.END, json.dumps(response_data, indent=4))  # отображаем ответ в поле ответа

# Создаем графическое окно
window = tk.Tk()

# Создаем поле для ввода URL
label_url = tk.Label(window, text="order_uid:")
label_url.pack()
entry_url = tk.Entry(window, width=25)
entry_url.pack()

# Создаем кнопку отправки запроса
button_send = tk.Button(window, text="Send", command=send_request)
button_send.pack()

# Создаем поле для вывода ответа
label_response = tk.Label(window, text="Response:")
label_response.pack()
entry_response = tk.Text(window,height=50)
entry_response.pack()

# Запускаем главный цикл окна
window.mainloop()