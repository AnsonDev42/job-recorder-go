# Job recorder go
Welcome to Job Recorder Go, an innovative solution designed to streamline the job application process for seekers everywhere. Born out of personal experience and the need for efficiency, this application automates the tedious task of tracking job applications. By leveraging OCR and LLM for summarization, Job Recorder Go transforms your job search recording experience into a manageable and insightful journey. 

## Motivation
The idea for Job Recorder Go sparked from my own challenges while navigating the sea of job applications. I found myself overwhelmed by the manual effort required to log and analyze each application. It was not only time-consuming but boring
I found myself overwhelmed by the manual effort required to log and analyze each application. It was not only time-consuming but also detracted from the time I could spend on actually applying or enhancing my skills. I hope this app can make life easier for everyone and the future me.

## A solution

Instead putting your just-applied job info into multiple columns in Notion, you can just screenshot the job descrption, click the upload, local prepressed OCR texts and let **GPT summarise**(WIP) the info such as job title (level), company and expected salary.

In addition, every day it sends you a email (for now only telegram bot notification) with the summarisation of your work today and encourage you ((WIP) if you reach the goal or not.

## Project usage
A single binary executable file. 

or compile/run by yourself: ```$git clone``` then ```$go run cmd/job-recorder-go/main.go```

## Project setup

This project strcture follows the recommended [Go project layout](https://github.com/golang-standards/project-layout)


## tech stack
This application is built with Golang and with [fyne](https://fyne.io/) for GUI framework.
This being my first Go project, I embraced the challenge of learning a new language and its ecosystem. I delved into Go's concurrency model, interface system, and package management, applying these concepts to develop this robust and efficient application.

## Road Map 
- ✔️ daily reminder (telegram bot, fill in GUI)
- ✔️ due to system limitation, currently only support MacOS's local OCR (otherwise need to use tesseract, a pain in the ass dynamtic lib link issue for installation in multi-platforms).
- 🚧 LLM support
