#!/usr/bin/env python3
'''
    File name: gui.py
    Author: Tomáš Vlček
    Date created: 2023-10-10
    License: GNU General Public License v3.0
    Python Version: ≥3.11.5
'''

import defopt
import PySimpleGUI as sg

def main():
    layout = [
        [sg.Graph((800, 800), (0, 450), (450, 0), key='-GRAPH-', enable_events=True, change_submits=True, drag_submits=False)],
        [
            sg.Input(key='-INPUT-'),
            sg.FileBrowse(file_types=(("PNG Images", "*.png"), ("ALL Files", "*.*"))),
            sg.Button('Load file', key='-OPEN_FILE-'),
            sg.Push(),
            sg.Button('Save file', key='-SAVE_FILE-')
        ],
    ]

    window = sg.Window('hans', layout)

    filename: str = str()

    while True:
        event, values = window.read()
        print(event, values)

        if event == sg.WINDOW_CLOSED or event == 'Quit':
            break
        elif event == '-OPEN_FILE-':
            filename = values['-INPUT-']
            print(f'Otvírám soubor {filename}')
        elif event == '-SAVE_FILE-':
            print(f'Ukládám soubor {filename}')
            # output_picture.picture.save_svg('output.svg')

     # Finish up by removing from the screen
    window.close()

    return

if __name__ == '__main__':
	defopt.run(main)