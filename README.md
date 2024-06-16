# Trade Republic Portfolio Downloader

## Preamble

This project exists because Trade Republic does not provide a comprehensive view of all purchases and sales. While it is possible to see current holdings in analytics, there is no way to view the analytics of sold assets and their profit. All sale transactions must be tracked manually to understand trading benefits better. This project was created to fill that gap by providing a better representation of trading activities.

## Existing Solutions

Research revealed a few solutions that fulfill similar requirements. However, many are outdated and do not utilize new endpoints (referred to here as websocket message types). Additionally, they are limited in tracking purchases and sales of assets.

Main disadvantages of existing solutions:

* Resetting the paired device (annoying if you use the mobile app regularly).
* Requiring OCR for transaction details (reading from transaction PDF documents).
* Using outdated endpoints to fetch data (no support for newly introduced "Save-back" and "Round up" transactions).
* Being written in Python (not our primary programming language).

## Implementation

### Minimum Requirements and Limitations

A few requirements and limitations were set initially and strictly followed during this project's planning and implementation:

* It should be written in Go and compiled as a binary for all major platforms and architectures.
* Make it open source to allow contributions and audits from others.
* No configuration file requirements (all input requested in console).
* No dependencies (such as SQL databases) to enable non-tech users to use the app "as is".
* Writing results into a CSV file to use it with Excel to build formulas, filter data, etc.
* No security information storage (except for session and refresh tokens) on the host machine.
* **No data should leave the host machine.**

The application performs the same functions as Trade Republic's official frontend application:

* Authenticates using the same API endpoints.
* Retrieves data using the same websocket address.

### Currently Supported Functionality

* Transaction PDF documents download.
* Account related PDF documents download (aka Aktivität).
* Creating a CSV file with all transactions (except for "interest received" transactions). This includes:
  * Cash deposits and withdrawals.
  * Purchase and sale of ETFs, stocks, and cryptocurrency.
  * Limited support for the purchase of derivatives.
  * Dividends received from ETFs and stocks.
  * Benefits received such as round-up and save-back.
* Inserting new data into the CSV file.
* Saving raw responses onto the file system.

### Planned Features and Improvements

**Upcoming Features:**

* Support for including "lending" transactions.
* Identifying stock transactions.
* Writing data into an SQLite file on the filesystem.
* Calculating miscellaneous values based on Trade Republic data: invested amount, taxable amount, earliest date of non-taxable sale of crypto assets, etc.
* Increasing source code test coverage.

**Potential Future Features:**

* Writing data into an SQL database for use in custom applications.
* Developing a frontend application to better visualize all transactions in a user-friendly way. More details will follow if development begins.

## Usage

### Choosing the Right Binary for Your OS and Architecture

Download one of the binaries from the [releases](https://github.com/dhojayev/traderepublic-portfolio-downloader/releases) section according to the table below:

| OS      | Architecture | Description                                            | File to download      |
| ------- | ------------ | ------------------------------------------------------ | --------------------- |
| macOS   | amd64        | Apple devices using Intel CPU                          | *-darwin-amd64.tar.gz |
| macOS   | arm64        | Apple devices using Apple Silicon (Apple M1 and newer) | *-darwin-arm64.tar.gz |
| Windows | amd64        | 64-Bit Windows                                         | *-windows-amd64.zip   |
| Windows | arm64        | Windows for ARM                                        | *-windows-arm64.zip   |
| Linux   | amd64        | 64-Bit Linux distro                                    | *-linux-amd64.tar.gz  |
| Linux   | arm64        | Linux distro for ARM processors                        | *-linux-arm64.tar.gz  |

*Users using Windows and macOS may receive a warning message before running the binary since it has not been signed.*

### Running the App

All available arguments and flags:

```shell
➜  traderepublic-portfolio-downloader git:(main) ✗ ./bin/portfoliodownloader --help
Usage: portfoliodownloader [--write-responses] [--debug] [--trace]

Options:
  --write-responses, -w
                         write API responses to the file system
  --debug                enable debug mode
  --trace                enable trace mode
  --help, -h             display this help and exit
```

After downloading a binary for your OS, run it in the terminal:

```shell
➜  traderepublic-portfolio-downloader git:(main) ✗ ./bin/portfoliodownloader
Enter phone number in international format (+49xxxxxxxxxxxxx): 
```

Enter your registered mobile number in international format:

```shell
➜  traderepublic-portfolio-downloader git:(main) ✗ ./bin/portfoliodownloader
Enter phone number in international format (+49xxxxxxxxxxxxx): 
+491234567890
```

Provide your PIN and hit enter:

```shell
➜  traderepublic-portfolio-downloader git:(main) ✗ ./bin/portfoliodownloader
Enter phone number in international format (+49xxxxxxxxxxxxx): 
+491234567890
Enter pin:
```

Enter the OTP received from Trade Republic and hit enter:

```shell
➜  traderepublic-portfolio-downloader git:(main) ✗ ./bin/portfoliodownloader
Enter phone number in international format (+49xxxxxxxxxxxxx): 
+491234567890
Enter pin:
Enter 2FA token:
```

You will see the progress of the download and processing:

```shell
➜  traderepublic-portfolio-downloader git:(main) ✗ ./bin/portfoliodownloader
Enter phone number in international format (+49xxxxxxxxxxxxx): 
+491234567890
Enter pin:
Enter 2FA token:
Mar 28 12:02:09.385 [INFO] 247 transactions downloaded
Mar 28 12:02:09.385 [INFO] [id:xxxxxxx-xxxxx-xxxxx-xxxxx-xxxx] Fetching transaction details
Mar 28 12:02:09.413 [INFO] [id:xxxxxxx-xxxxx-xxxxx-xxxxx-xxxx] Processing transaction details
Mar 28 12:02:09.453 [INFO] [id:xxxxxxx-xxxxx-xxxxx-xxxxx-xxxx] Unsupported transaction skipped
Mar 28 12:02:09.453 [INFO] [id:xxxxxxx-xxxxx-xxxxx-xxxxx-xxxx] Fetching transaction details
Mar 28 12:02:09.485 [INFO] [id:xxxxxxx-xxxxx-xxxxx-xxxxx-xxxx] Processing transaction details
Mar 28 12:02:09.488 [INFO] [id:xxxxxxx-xxxxx-xxxxx-xxxxx-xxxx] Unsupported transaction skipped
...
Mar 28 12:02:27.379 [INFO] Transactions completed: 200, skipped: 47
Jun 16 22:15:09.675 [INFO] Downloading activity log entries
Jun 16 22:15:09.858 [INFO] 42 activity log entries downloaded
Jun 16 22:15:09.858 [INFO] [id:14e2842f-7a4a-4f3f-8f37-52ac2d332673] Fetching activity log entry details
Jun 16 22:15:10.030 [INFO] [id:14e2842f-7a4a-4f3f-8f37-52ac2d332673] Document downloaded
Jun 16 22:15:10.031 [INFO] [id:f61e9ebb-3e6b-473b-91ae-1b61bc8b9dcc] Fetching activity log entry details
Jun 16 22:15:10.133 [INFO] [id:f61e9ebb-3e6b-473b-91ae-1b61bc8b9dcc] Document downloaded
...
Jun 16 22:15:13.484 [INFO] Activity log entries completed: 29, skipped: 13
```

### PDF documents storage structure

All documents are saved under `documents` directory followed by:
* `transactions` or `activity` based on document type.
* Date in `YYYY-mm` format, e.g.: `2020-01`
* UUID related to its transaction or activity log entry. 
* Document name, e.g. `Abrechnung Ausführung.pdf` 

Exmaples:
* `documents/transactions/2023-11/67bd11bb-327e-475d-b715-5876bab61c5c/Abrechnung Ausführung.pdf`
* `documents/activity/2023-09/d7488784-2ac1-4daf-856a-f87fb79f641e/Kundenvereinbarung.pdf`

### CSV File Fields

| Field              | Description                                                               |
| ------------------ | ------------------------------------------------------------------------- |
| **ID**             | Transaction UUID                                                          |
| **Status**         | Transaction status (should always be `executed`)                          |
| **Timestamp**      | Date and time of transaction execution, e.g.: `30 Nov 23 10:22 +0000`     |
| **Type**           | Transaction type, one of: `Purchase, Sale, Dividends, Round Up, Saveback` |
| **Asset type**     | Asset type, one of: `ETF, Cryptocurrency, Lending, Other`                 |
| **Name**           | Asset name, e.g.: `Bitcoin`                                               |
| **Instrument**     | Instrument ISIN, e.g.: `IE00BK1PV551`                                     |
| **Shares**         | Number of shares in the transaction (negative when sold)                  |
| **Rate**           | Price per share in EUR                                                    |
| **Realized yield** | Realized yield in percentage (negative if loss)                           |
| **Realized PnL**   | Realized profit or loss amount in EUR (negative if loss)                  |
| **Commission**     | Commission paid to Trade Republic for the transaction in EUR              |
| **Debit**          | Amount debited from the deposited amount in EUR                           |
| **Credit**         | Amount credited to the deposited amount in EUR                            |
| **Tax amount**     | Tax applied to this transaction in EUR                                    |
| **Documents**      | File path to the document PDF files                                       |

Example CSV output can be viewed here: [transactions.csv](./assets/transactions.csv)

**All values saved into the CSV file are taken "as is" from Trade Republic (except for making them negative in respective cases).**

## Troubleshooting

You may encounter errors while running the app due to unexpected data. This is normal as we could only cover the data in our portfolios. It is possible that you have some types of assets we do not own and cannot test, such as derivatives.

Please create an issue and attach the failing response with falsified amounts and removed ID. We will either implement support for these responses or ensure the app does not fail when such a response is received.

## Have Suggestions or Improvements?

We hope this app will improve with community help until Trade Republic implements a better dashboard for an overview of such data.

Please create a pull request with your changes if you have something to contribute. We are very open to constructive suggestions and feedback.

## Closing Words

This project and its contributors have no affiliation with Trade Republic Bank GmbH.

Trade Republic is a registered trademark of Trade Republic Bank GmbH.