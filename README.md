# Trade Republic Portfolio Downloader

## Preamble

This project exists solely because I need a better representation of all purchases and sales made in Trade Republic.

The initial idea was to create a list of all purchase and sale transactions to be able to track profit and loss.
Unfortunately Trade Republic does not provide such dashboard or file(s) to download (it is possible to view current
portfolio
assets but no analytics of what has been sold and how much of profit has been gained).

## Existing solutions

After researching a bit in internet I found a solution (or
two) that fulfill similar requirements, however, those
solutions are pretty outdated and do not benefit from new endpoints (actually websocket message types, but let's call
them here that way for simplicity).

Main disadvantages of existing solutions for me:

* Resets paired device (which is annoying if you use the mobile app regularly)
* Requires OCR for getting transaction details
* Pretty outdated and not well-maintained
* Written in python (which I am not proficient in enough to contribute)

## Implementation

### Minimum requirements and limitations

I initially set a few requirements and limitations for myself when planning this project which I followed:

* Written in go and compiled as a binary for all major platforms and architectures
* Making it opensource to allow others to contribute and audit
* No configuration file requirements (all input requested in console)
* No dependencies (such as SQL databases, etc) to enable non-techs using the app "AS IS"
* Writing results into a CSV file to be able to build formulas, filter, etc. the data
* No security information storage (except for session and refresh tokens) on host machine
* **No data should leave the host machine**

It is important to mention that this application does nothing more than Trade Republic's official frontend application
would do:

* Authenticates using the same API endpoints
* Retrieves the data using the same websocket address

### Currently supported functionality

* Creating CSV file with all transaction (except for "interest received" transaction for now). This includes:
  * purchase and sale of ETFs, stocks, cryptocurrency;
  * dividends received from ETFs and stocks;
  * benefits received suchs as round up and save-back;
* Inserting new data into the CSV file
* Saving raw responses onto the file system

### Planned features and improvements

**What is coming:**

* Support of "interest received" transactions
* Identifying stock transactions
* Downloading and storing PDF files attached to each transaction
* Source code test coverage

**What (maybe) will follow:**

* Writing the data into an SQL database for using it in custom applications

## Usage

### Choosing binary for your OS and architecture:

Download one of the binaries from [releases](https://github.com/dhojayev/traderepublic-portfolio-downloader/releases)
section according to the table below:

| OS      | Architecture | Description                          | File to download      |
|---------|--------------|--------------------------------------|-----------------------|
| macOS   | amd64        | Apple devices using Intel CPU        | *-darwin-amd64.tar.gz |
| macOS   | arm64        | Apple devices using M1 SoC and newer | *-darwin-arm64.tar.gz |
| Windows | amd64        | 64-Bit Windows                       | *-windows-amd64.zip   |
| Linux   | amd64        | 64-Bit linux distro                  | *-linux-amd64.tar.gz  |
| Linux   | arm64        | linux distro for ARM processors      | *-linux-arm64.tar.gz  |


*Users using macOS may require to allow running the binary since it has not been signed.*

### Running the app

All available arguments and flags:

```shell
➜  traderepublic-portfolio-downloader git:(main) ✗ ./bin/portfoliodownloader --help
Usage: portfoliodownloader [--write-responses] [--debug] [--trace] PHONENUMBER

Positional arguments:
  PHONENUMBER            Phone number in international format: +49xxxxxxxxxxxxx

Options:
  --write-responses, -w
                         Write api responses to file system
  --debug                Enable debug mode
  --trace                Enable trace mode
  --help, -h             display this help and exit

```

After downloading a binary for respective OS simply run it in terminal by providing your phone number as an argument:

```shell
➜  traderepublic-portfolio-downloader git:(main) ✗ ./bin/portfoliodownloader +49xxxxxxxxxxxx
Enter pin:
```

Provide your pin and hit enter:

```shell
➜  traderepublic-portfolio-downloader git:(main) ✗ ./bin/portfoliodownloader +49xxxxxxxxxxxx
Enter pin:
Mar 28 12:01:32.249 [INFO] Downloading transactions
Enter 2FA token:
```

Enter OTP that you received from Trade Republic and hit enter:

```shell
➜  traderepublic-portfolio-downloader git:(main) ✗ ./bin/portfoliodownloader +49xxxxxxxxxxxx
Enter pin:
Mar 28 12:01:32.249 [INFO] Downloading transactions
Enter 2FA token:
1111
```

You will see the progress of download and processing:

```shell
➜  traderepublic-portfolio-downloader git:(main) ✗ ./bin/portfoliodownloader +49xxxxxxxxxxxx
Enter pin:
Mar 28 12:01:32.249 [INFO] Downloading transactions
Enter 2FA token:
1111
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

| Field               | Description                                                               |
|---------------------|---------------------------------------------------------------------------|
| **ID**              | Transaction UUID                                                          |
| **Status**          | Transaction status (should always be `executed`)                          |
| **Timestamp**       | Date and time of transaction execution, e.g.: `30 Nov 23 10:22 +0000`     |
| **Type**            | Transaction type, one of: `Purchase, Sale, Dividends, Round Up, Saveback` |
| **Asset type**      | Asset type, one of: `ETF, Cryptocurrency, Lending, Other`                 |
| **Name**            | Asset name, e.g.: `Bitcoin`                                               |
| **Instrument**      | Instrument ISIN, e.g.: `IE00BK1PV551`                                     |
| **Shares**          | Number of shares in transaction                                           |
| **Rate**            | Price per share in EUR                                                    |
| **Realized yield**  | Realized yield in percentage                                              |
| **Realized PnL**    | Realized profit or loss amount in EUR (negative is loss)                  |
| **Commission**      | Commission paid to Trade Republic for the transaction in EUR              |
| **Debit**           | Amount debited from the deposited amount in EUR                           |
| **Credit**          | Amount credited to the deposited amount in EUR                            |
| **Portfolio value** | Amount that contributes to the portfolio size in EUR                      |

Example CSV output can be viewed here: [transactions.csv](./assets/transactions.csv)

## Troubleshooting

It is possible that you will get an error while running the app because it will receive unexpected data. This is normal
because I was able to cover only the data I have in my portfolio. It is possible that you have some type of assets that
I simply don't own and cannot test, e.g: derivates.

Please create an issue and attach failing response with falsified values in it. I will the either add support of these
responses or make sure to not let the app fail when received such response.

## Have suggestions or improvements?

I hope that this app will become better with the help of the community until Trade Republic decides to implement a
better dashboard to be able to have an overview of such data.

Please create a pull request with your changes if you have something to contribute. I am very open for constructive
suggestions and feedback.

## Closing words

This project and I have no affiliation to Trade Republic Bank GmbH by any means.

Trade Republic is a registered trademark of Trade Republic Bank GmbH.
