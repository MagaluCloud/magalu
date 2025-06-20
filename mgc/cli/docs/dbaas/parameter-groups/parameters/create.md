# Create

Create a parameter for a group.

## Usage:
```
mgc dbaas parameter-groups parameters create [parameter-group-id] [flags]
```

## Examples:
```
mgc dbaas parameter-groups parameters create --name="LOWER_CASE_TABLE_NAMES"
```

## Flags:
```
    --cli.list-links enum[=table]   List all available links for this command (one of "json", "table" or "yaml")
-h, --help                          help for create
    --name enum                     EngineParametersMysql84 (one of "ACTIVATE_ALL_ROLES_ON_LOGIN", "AUTOCOMMIT", "AUTOMATIC_SP_PRIVILEGES", "BACK_LOG", "BLOCK_ENCRYPTION_MODE", "CHARACTER_SET_SERVER", "COLLATION_SERVER", "CONNECT_TIMEOUT", "CTE_MAX_RECURSION_DEPTH", "DEFAULT_WEEK_FORMAT", "DIV_PRECISION_INCREMENT", "END_MARKERS_IN_JSON", "EQ_RANGE_INDEX_DIVE_LIMIT", "EVENT_SCHEDULER", "FT_QUERY_EXPANSION_LIMIT", "GENERAL_LOG", "GENERATED_RANDOM_PASSWORD_LENGTH", "GROUP_CONCAT_MAX_LEN", "INNODB_BUFFER_POOL_INSTANCES", "INNODB_FILL_FACTOR", "INNODB_FT_ENABLE_STOPWORD", "INNODB_FT_MAX_TOKEN_SIZE", "INNODB_FT_MIN_TOKEN_SIZE", "INNODB_FT_NUM_WORD_OPTIMIZE", "INNODB_FT_RESULT_CACHE_LIMIT", "INNODB_LOCK_WAIT_TIMEOUT", "INNODB_LOG_BUFFER_SIZE", "INNODB_MAX_UNDO_LOG_SIZE", "INNODB_OPTIMIZE_FULLTEXT_ONLY", "INNODB_PRINT_ALL_DEADLOCKS", "INNODB_PRINT_DDL_LOGS", "INNODB_READ_IO_THREADS", "INNODB_REDO_LOG_CAPACITY", "INNODB_ROLLBACK_ON_TIMEOUT", "INNODB_SEGMENT_RESERVE_FACTOR", "INNODB_WRITE_IO_THREADS", "INTERACTIVE_TIMEOUT", "JOIN_BUFFER_SIZE", "LC_TIME_NAMES", "LOCK_WAIT_TIMEOUT", "LOG_ERROR_VERBOSITY", "LOG_QUERIES_NOT_USING_INDEXES", "LOG_SLOW_ADMIN_STATEMENTS", "LOG_THROTTLE_QUERIES_NOT_USING_INDEXES", "LONG_QUERY_TIME", "LOWER_CASE_TABLE_NAMES", "MAX_ALLOWED_PACKET", "MAX_CONNECTIONS", "MAX_CONNECT_ERRORS", "MAX_EXECUTION_TIME", "MAX_JOIN_SIZE", "MAX_POINTS_IN_GEOMETRY", "MAX_PREPARED_STMT_COUNT", "MAX_SEEKS_FOR_KEY", "MAX_SORT_LENGTH", "MAX_SP_RECURSION_DEPTH", "MAX_USER_CONNECTIONS", "MAX_WRITE_LOCK_COUNT", "MIN_EXAMINED_ROW_LIMIT", "NET_READ_TIMEOUT", "NET_RETRY_COUNT", "NET_WRITE_TIMEOUT", "NGRAM_TOKEN_SIZE", "PARSER_MAX_MEM_SIZE", "PASSWORD_REQUIRE_CURRENT", "PASSWORD_REUSE_INTERVAL", "QUERY_ALLOC_BLOCK_SIZE", "REGEXP_STACK_LIMIT", "REGEXP_TIME_LIMIT", "SHOW_CREATE_TABLE_VERBOSITY", "SKIP_NAME_RESOLVE", "SKIP_SHOW_DATABASE", "SLOW_LAUNCH_TIME", "SLOW_QUERY_LOG", "SQL_MODE", "STORED_PROGRAM_CACHE", "STORED_PROGRAM_DEFINITION_CACHE", "TABLESPACE_DEFINITION_CACHE", "TABLE_DEFINITION_CACHE", "TABLE_OPEN_CACHE", "TABLE_OPEN_CACHE_INSTANCES", "TEMPTABLE_MAX_MMAP", "TEMPTABLE_MAX_RAM", "TMP_TABLE_SIZE", "TRANSACTION_ISOLATION", "UPDATABLE_VIEWS_WITH_LIMIT", "WAIT_TIMEOUT" or "WINDOWING_USE_HIGH_PRECISION") (required)
    --parameter-group-id uuid       Value referring to parameter group Id. (required)
    --value anyOf                   Value (at least one of: number, integer, boolean or string)
                                    Use --value=help for more details (required)
```

## Global Flags:
```
    --api-key string           Use your API key to authenticate with the API
-U, --cli.retry-until string   Retry the action with the same parameters until the given condition is met. The flag parameters
                               use the format: 'retries,interval,condition', where 'retries' is a positive integer, 'interval' is
                               a duration (ex: 2s) and 'condition' is a 'engine=value' pair such as "jsonpath=expression"
-t, --cli.timeout duration     If > 0, it's the timeout for the action execution. It's specified as numbers and unit suffix.
                               Valid unit suffixes: ns, us, ms, s, m and h. Examples: 300ms, 1m30s
    --debug                    Display detailed log information at the debug level
    --env enum                 Environment to use (one of "pre-prod" or "prod") (default "prod")
    --no-confirm               Bypasses confirmation step for commands that ask a confirmation from the user
-o, --output string            Change the output format. Use '--output=help' to know more details.
-r, --raw                      Output raw data, without any formatting or coloring
    --region enum              Region to reach the service (one of "br-mgl1", "br-ne1" or "br-se1") (default "br-se1")
    --server-url uri           Manually specify the server to use
```

