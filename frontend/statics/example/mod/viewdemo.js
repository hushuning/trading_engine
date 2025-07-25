layui.define(['form',"baseinfo", 'utils', 'kchart', 'websocket','login'], function(exports){
    var baseinfo = layui.baseinfo;
    var login = layui.login;

    var layer = layui.layer //弹层
        , form = layui.form
        , utils = layui.utils
        , $ = layui.$;
        

    var obj = {
        bind: function(){
            var me = this;
            form.on('select(order_type)', function (data) {
                if (data.value == "limit") {
                    $(".item-price").show();
                    $(".item-quantity").show();
                    $(".item-amount").hide();
                    $(".item-market-type").hide();
                    $(".qty-tips").hide();
                } else if (data.value == "market") {
                    $(".item-price").hide();
                    $(".item-market-type").show();
                    $(".qty-tips").show();
                }
                form.render('select');
            });
            form.on('radio(market-type)', function (data) {
                if (data.value == "market_qty") {
                    $(".item-quantity").show();
                    $(".item-amount").hide();
                    $(".qty-tips").show();
                } else {
                    $(".item-quantity").hide();
                    $(".qty-tips").hide();
                    $(".item-amount").show();
                }
            });

            //取消订单
            $("body").on("click", ".cancel-order", function(){
                var order_id = $(this).parents("tr").attr("order-id");
                console.log(order_id);

                $.ajax({
                    url: "/api/v1/base/order/cancel",
                    type: "post",
                    dataType: "json",
                    contentType: "application/json",
                    beforeSend: function(r) {
                        r.setRequestHeader("token", Cookies.get("user_id"));
                    },
                    data: function () {
                        var data = {
                            symbol: CURRENT_SYMBOL,
                            order_id: order_id,
                        };
                        return JSON.stringify(data)
                    }(),
                    success: function (d) {
                        if(d.ok){
                            layer.msg("已提交")
                        }else{
                            layer.msg(d.reason);
                        }
                    }
                });
            }).on("click", ".get_original_assets", function(){
                $.ajax({
                    url: "/example/deposit",
                    type: "get",
                    data: {
                        symbol: baseinfo.cfg_info.target.symbol + "," + baseinfo.cfg_info.base.symbol
                    },
                    contentType: "application/json",
                    success: function (d) {
                        if(d.ok){
                            me.load_assets();
                        }else{
                            layer.msg(d.reason);
                        }
                    }
                });
            });

            //新订单
            $(".opt").on("click", function () {
                var side = $(this).hasClass("sell") ? "ask" : "bid";
                var order_type = $("select[name='order_type']").val();
                var mtype = $("input[name='mtype']:checked").val();
                console.log("side: ", side, " order_type: ", order_type, " mtype: ", mtype);
                $.ajax({
                    url: "/api/v1/order/create",
                    type: "post",
                    dataType: "json",
                    contentType: "application/json",
                    data: function () {
                        var data = {
                            symbol: CURRENT_SYMBOL,
                            side: side,
                            order_type: order_type,
                        };

                        if (order_type == "market") {
                            if (mtype == "market_qty") {
                                data.qty = $("input[name='quantity']").val();
                            } else {
                                data.amount = $("input[name='amount']").val();
                            }
                        } else {
                            data.price = $("input[name='price']").val();
                            data.qty = $("input[name='quantity']").val();
                        }

                        console.log(data);
                        return JSON.stringify(data)
                    }(),
                    success: function (d) {
                        if(d.ok){
                            layer.msg("下单成功");
                            
                            me.load_assets();
                            me.load_order_unfinished();
                        }else{
                            layer.msg(d.reason);
                        }
                    }
                });
            });
        },

        load_depth_data: function(){
            $.get("/api/v1/market/depth?symbol="+CURRENT_SYMBOL+"&limit=10", function(d){
                if(d.ok){
                    utils.renderdepth(d.data);
                }
            });
        },
        load_tradelog_data: function(){
            $.get("/api/v1/market/trades?symbol="+CURRENT_SYMBOL+"&limit=10", function (d) {
                if (d.ok) {
                    var data = d.data.reverse();
                    for(var i=0; i<data.length; i++){
                        utils.rendertradelog(data[i]);
                    }

                }
            });
        },
        load_system_info: function(){
            $.get("/api/v1/quote/system", function(d){
                $(".version").html(d.version);
                $(".build").html(d.build);
            });
        },
        load_assets: function(){
            console.log('info:', baseinfo.cfg_info);
            $.ajax({
                url: "/api/v1/user/asset/query",
                type: "get",
                data:{
                    symbols: baseinfo.cfg_info.target.symbol+ "," + baseinfo.cfg_info.base.symbol,
                    t: Date.now()
                },
                success: function (d) {
                    console.log("load_assets: ", d);
                    if(!d.ok) {
                        layer.msg(d.reason);
                        return;
                    }

                    var html = [];
                    if(d.data.length > 0){
                        for(var i=0; i<d.data.length; i++){
                            html.push(" " + d.data[i].symbol.toUpperCase() + ":" + d.data[i].avail_balance);
                        }
                    }else{
                        html.push("<a href='javascript:;' style='color:red;' class='get_original_assets'>点我获取资产</a>");
                    }
                    console.log(html);
                    $(".assets .list").html(html.join(" "));
                }
            });
        },
        load_order_unfinished: function(){
            $.ajax({
                url: "/api/v1/order/unfinished",
                type: "get",
                data:{
                    symbol: baseinfo.cfg_info.symbol,
                    limit: 4,
                    t: Date.now()
                },
                success: function (d) {
                    console.log("load_order_unfinished: ", d);
                    if(d.ok){
                        $(".myorder-item").remove();
                        if(d.data.length > 0){
                            var data = d.data.reverse();
                            for(var i=0; i<data.length; i++){
                                utils.rendermyorder(data[i]);
                            }
                        }
                    }
                }
            });
        },
        load_all_tsymbols: function(){
            $.ajax({
                url: "/api/v1/base/trading/varieties",
                type: "get",
                success: function (d) {
                    console.log("/trading/varieties: ", d);
                    if(d.ok){
                        var data = d.data;
                        if(data.length > 0){
                            var html = [];
                            for(var i=0; i<data.length; i++){
                                html.push('<a href="/'+ data[i].symbol +'"><b>'+ data[i].symbol.toUpperCase() +'</b></a>');
                            }
                            $(".header-all-symbols").html(html.join(""))
                        }
                    }
                }
            });
        },
        init: function(){
            login.init();
            this.bind();
            // this.load_system_info();
            // this.load_all_tsymbols();
            this.load_depth_data();
            this.load_tradelog_data();
            this.load_assets();
            this.load_order_unfinished();
            // websocket.init();
        }
    };
    
    obj.init();
    exports('viewdemo', obj);
});

