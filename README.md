# Trade Republic Portfolio Downloader

## Preamble

This project exists because Trade Republic does not provide a better representation of all purchases and sales. While it is still possible to see current holdings in analytics (it is possible to view current
portfolio assets but no analytics of what has been sold and how much of profit has been gained), all sale transactions have to be tracked manually one-by-one to understand your benefits of trading better. Since Trade Republic does not provide such dashboard or file(s) to download this project's idea was born.

## Existing solutions

After researching couple of solutions that fulfill similar requirements were discovered, however, some
are pretty outdated and do not benefit from new endpoints (actually websocket message types, but let's call
them here that way for simplicity), some are very limited in fullfiling the main requirements: tracking purchases and sales of assets.

Main disadvantages of existing solutions in our opinion were:

* Reseting paired device (which is annoying if you use the mobile app regularly);
* Requiring OCR for getting transaction details (reading from transaction PDF documents);
* Using outdated endpoints to fetch data (no support of newly introduced "Save-back" and "Round up" transactions);
* Being written in python (which is not the main language of programming of ours);

## Implementation

### Minimum requirements and limitations

A few requirements and limitations were set intially when planning this project which we strictly follow:

* It should be written in go and compiled as a binary for all major platforms and architectures;
* Making it opensource to allow others to contribute and audit;
* No configuration file requirements (all input requested in console);
* No dependencies (such as SQL databases, etc) to enable non-techs using the app "AS IS";
* Writing results into a CSV file to be able to build formulas, filter the data, etc.;
* No security information storage (except for session and refresh tokens) on host machine;
* **No data should leave the host machine;**

It is important to understand that this application does nothing more than Trade Republic's official frontend application would do:

* Authenticates using the same API endpoints;
* Retrieves the data using the same websocket address;

### Currently supported functionality

* Creating CSV file with all transaction (except for "interest received" transaction for now). This includes:
  * deposit transactions;
  * purchase and sale of ETFs, stocks, cryptocurrency;
  * interest received transactions;
  * limited support of purchase of derivatives;
  * dividends received from ETFs and stocks;
  * benefits received suchs as round up and save-back;
* Inserting new data into the CSV file
* Saving raw responses onto the file system

### Planned features and improvements

**What is coming:**

* Support of including "lending" transactions;
* Downloading and storing PDF files attached to each transaction;
* Identifying stock transactions;
* Writing data into an sqlite file on the filesystem;
* Calculating miscelaneous values based on data from TR: invested amount, taxable amount, earliest date of non-taxable sale of Crypto assets, etc;
* Source code test coverage;

**What (maybe) will follow:**

* Writing data into an SQL database for using it in custom applications;
* Frontend application to better visualize all transactions in a user-friendly way. More details will follow once (or if) the development starts;

## Usage

### Choosing the right binary for your OS and architecture:

Download one of the binaries from [releases](https://github.com/dhojayev/traderepublic-portfolio-downloader/releases)
section according to the table below:

| OS      | Architecture | Description                                            | File to download      |
| ------- | ------------ | ------------------------------------------------------ | --------------------- |
| macOS   | amd64        | Apple devices using Intel CPU                          | *-darwin-amd64.tar.gz |
| macOS   | arm64        | Apple devices using Apple Silicon (Apple M1 and newer) | *-darwin-arm64.tar.gz |
| Windows | amd64        | 64-Bit Windows                                         | *-windows-amd64.zip   |
| Windows | arm64        | Windows for ARM                                        | *-windows-arm64.zip   |
| Linux   | amd64        | 64-Bit linux distro                                    | *-linux-amd64.tar.gz  |
| Linux   | arm64        | linux distro for ARM processors                        | *-linux-arm64.tar.gz  |


*Users using Windows and macOS may get a warning message before running the binary since it has not been signed.*

### Running the app

All available arguments and flags:

```shell
➜  traderepublic-portfolio-downloader git:(main) ✗ ./bin/portfoliodownloader --help
Usage: portfoliodownloader [--write-responses] [--debug] [--trace]

Options:
  --write-responses, -w
                         Write api responses to file system
  --debug                Enable debug mode
  --trace                Enable trace mode
  --help, -h             display this help and exit

```

After downloading a binary for respective OS simply run it in terminal:

```shell
➜  traderepublic-portfolio-downloader git:(main) ✗ ./bin/portfoliodownloader
Enter phone number in international format (+49xxxxxxxxxxxxx): 
```

Enter your registered mobile number in international format as requested:

```shell
➜  traderepublic-portfolio-downloader git:(main) ✗ ./bin/portfoliodownloader
Enter phone number in international format (+49xxxxxxxxxxxxx): 
+491234567890
```

Provide your pin and hit enter:

```shell
➜  traderepublic-portfolio-downloader git:(main) ✗ ./bin/portfoliodownloader
Enter phone number in international format (+49xxxxxxxxxxxxx): 
+491234567890
Enter pin:
```

Enter OTP that you received from Trade Republic and hit enter:

```shell
➜  traderepublic-portfolio-downloader git:(main) ✗ ./bin/portfoliodownloader
Enter phone number in international format (+49xxxxxxxxxxxxx): 
+491234567890
Enter pin:
Enter 2FA token:
```

You will see the progress of download and processing:

```shell
➜  traderepublic-portfolio-downloader git:(main) ✗ ./bin/portfoliodownloader
Enter phone number in international format (+49xxxxxxxxxxxxx): 
+491234567890
Enter pin:
Enter 2FA token:
Mar 28 12:02:09.385 [INFO] 247 transaction downloaded
Mar 28 12:02:09.385 [INFO] [id:xxxxxxx-xxxxx-xxxxx-xxxxx-xxxx] Fetching transaction details
Mar 28 12:02:09.413 [INFO] [id:xxxxxxx-xxxxx-xxxxx-xxxxx-xxxx] Processing transaction details
Mar 28 12:02:09.453 [INFO] [id:xxxxxxx-xxxxx-xxxxx-xxxxx-xxxx] Unsupported transaction skipped
Mar 28 12:02:09.453 [INFO] [id:xxxxxxx-xxxxx-xxxxx-xxxxx-xxxx] Fetching transaction details
Mar 28 12:02:09.485 [INFO] [id:xxxxxxx-xxxxx-xxxxx-xxxxx-xxxx] Processing transaction details
Mar 28 12:02:09.488 [INFO] [id:xxxxxxx-xxxxx-xxxxx-xxxxx-xxxx] Unsupported transaction skipped
...
Mar 28 12:02:27.379 [INFO] All data processed
```

### CSV file fields

| Field              | Description                                                               |
| ------------------ | ------------------------------------------------------------------------- |
| **ID**             | Transaction UUID                                                          |
| **Status**         | Transaction status (should always be `executed`)                          |
| **Timestamp**      | Date and time of transaction execution, e.g.: `30 Nov 23 10:22 +0000`     |
| **Type**           | Transaction type, one of: `Purchase, Sale, Dividends, Round Up, Saveback` |
| **Asset type**     | Asset type, one of: `ETF, Cryptocurrency, Lending, Other`                 |
| **Name**           | Asset name, e.g.: `Bitcoin`                                               |
| **Instrument**     | Instrument ISIN, e.g.: `IE00BK1PV551`                                     |
| **Shares**         | Number of shares in transaction (negative when sold)                      |
| **Rate**           | Price per share in EUR                                                    |
| **Realized yield** | Realized yield in percentage (negative if loss)                           |
| **Realized PnL**   | Realized profit or loss amount in EUR (negative if loss)                  |
| **Commission**     | Commission paid to Trade Republic for the transaction in EUR              |
| **Debit**          | Amount debited from the deposited amount in EUR                           |
| **Credit**         | Amount credited to the deposited amount in EUR                            |

Example CSV output can be viewed here: [transactions.csv](./assets/transactions.csv)

**All values saved into the CSV file are taken "as is" from Trade Republic (except for making them negative in respective cases).**

## Troubleshooting

It is possible that you will get an error while running the app because it will receive unexpected data. This is normal because we were able to cover only the data we have in our portfolios. It is possible that you have some type of assets that we simply don't own and cannot test, e.g: derivates.

Please create an issue and attach failing response with falsified amounts and removed ID. We will the either implement support of these responses or make sure to not let the app fail when received such response.

## Have suggestions or improvements?

We hope that this app will become better with the help of the community until Trade Republic decides to implement a
better dashboard to be able to have an overview of such data.

Please create a pull request with your changes if you have something to contribute. We are very open for constructive
suggestions and feedback.

## Closing words

This project and it's contributors have no affiliation to Trade Republic Bank GmbH by any means.

Trade Republic is a registered trademark of Trade Republic Bank GmbH.
