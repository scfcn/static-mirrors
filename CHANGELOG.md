# 1.0.0 (2026-02-01)


### Bug Fixes

* 减少服务重启和健康检查的重试次数 ([d15e702](https://github.com/scfcn/static-mirrors/commit/d15e7024992fa7a57f1c03eb258fbacc6a0ee5d8))


### Features

* **proxy:** 增强缓存策略和路径代理功能 ([733c725](https://github.com/scfcn/static-mirrors/commit/733c725c0ff51d2c02b86e94eefed839082b54fe))
* 初始化前端文件公益镜像服务项目 ([dbe6dfd](https://github.com/scfcn/static-mirrors/commit/dbe6dfdbf5a30edd075b600e29e79bd023a35ac3))
* **缓存/统计:** 添加Redis支持并重构相关模块 ([c67234a](https://github.com/scfcn/static-mirrors/commit/c67234a58f033c5090579b7ac0c19b37ca1d0865))

# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [1.0.0] - 2026-02-01
### Added
- Initial release of static-mirrors
- Support for multiple sources: jsdelivr, cdnjs, ghcr, docker, unpkg
- Memory cache implementation
- SQLite statistics storage
- Docker deployment support
